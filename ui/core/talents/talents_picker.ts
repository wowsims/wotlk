import { Component } from '/tbc/core/components/component.js';
import { Spec } from '/tbc/core/proto/common.js';
import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { SpecTalents } from '/tbc/core/proto_utils/utils.js';
import { Player } from '/tbc/core/player.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { isRightClick } from '/tbc/core/utils.js';
import { sum } from '/tbc/core/utils.js';

const MAX_TALENT_POINTS = 61;
const NUM_ROWS = 9;
const TALENTS_STORAGE_KEY = 'Talents';

export abstract class TalentsPicker<SpecType extends Spec> extends Component {
	private readonly player: Player<SpecType>;
	frozen: boolean;
	readonly trees: Array<TalentTreePicker<SpecType>>;

	constructor(parent: HTMLElement, player: Player<SpecType>, treeConfigs: TalentsConfig<SpecType>) {
		super(parent, 'talents-picker-root');
		this.player = player;
		this.frozen = false;
		this.trees = treeConfigs.map(treeConfig => new TalentTreePicker(this.rootElem, player, treeConfig, this));
		this.trees.forEach(tree => tree.talents.forEach(talent => talent.setPoints(0, false)));

		this.setTalentsString(TypedEvent.nextEventID(), this.player.getTalentsString());
		this.player.talentsChangeEmitter.on(eventID => {
			this.setTalentsString(eventID, this.player.getTalentsString());
		});
	}

	get numPoints() {
		return sum(this.trees.map(tree => tree.numPoints));
	}

	isFull() {
		return this.numPoints >= MAX_TALENT_POINTS;
	}

	update(eventID: EventID) {
		if (this.isFull()) {
			this.rootElem.classList.add('talents-full');
		} else {
			this.rootElem.classList.remove('talents-full');
		}

		this.trees.forEach(tree => tree.update());

		TypedEvent.freezeAllAndDo(() => {
			this.player.setTalentsString(eventID, this.getTalentsString());
		});
	}

	getTalentsString(): string {
		return this.trees.map(tree => tree.getTalentsString()).join('-').replace(/-+$/g, '');
	}

	setTalentsString(eventID: EventID, str: string) {
		const parts = str.split('-');
		this.trees.forEach((tree, idx) => tree.setTalentsString(parts[idx] || ''));
		this.update(eventID);
	}

	// Freezes the talent calculator so that user input cannot change it.
	freeze() {
		this.frozen = true;
		this.rootElem.classList.add('frozen');
	}
}

class TalentTreePicker<SpecType extends Spec> extends Component {
	private readonly config: TalentTreeConfig<SpecType>;
	private readonly title: HTMLElement;

	readonly talents: Array<TalentPicker<SpecType>>;
	readonly picker: TalentsPicker<SpecType>;

	// The current number of points in this tree
	numPoints: number;

	constructor(parent: HTMLElement, player: Player<SpecType>, config: TalentTreeConfig<SpecType>, picker: TalentsPicker<SpecType>) {
		super(parent, 'talent-tree-picker-root');
		this.config = config;
		this.numPoints = 0;
		this.picker = picker;

		this.rootElem.innerHTML = `
    <div class="talent-tree-header">
      <span class="talent-tree-title"></span>
      <span class="talent-tree-reset fa fa-times"></span>
    </div>
    <div class="talent-tree-main">
    </div>
    `;

		this.title = this.rootElem.getElementsByClassName('talent-tree-title')[0] as HTMLElement;

		const main = this.rootElem.getElementsByClassName('talent-tree-main')[0] as HTMLElement;
		main.style.backgroundImage = `url('${config.backgroundUrl}')`;

		this.talents = config.talents.map(talent => new TalentPicker(main, player, talent, this));
		this.talents.forEach(talent => {
			if (talent.config.prereqLocation) {
				this.getTalent(talent.config.prereqLocation).config.prereqOfLocation = talent.config.location;
			}
		});

		const reset = this.rootElem.getElementsByClassName('talent-tree-reset')[0] as HTMLElement;
		reset.addEventListener('click', event => {
			if (!this.picker.frozen) {
				this.talents.forEach(talent => talent.setPoints(0, false));
				this.picker.update(TypedEvent.nextEventID());
			}
		});
	}

	update() {
		this.title.textContent = this.config.name + ' (' + this.numPoints + ')';
		this.talents.forEach(talent => talent.update());
	}

	getTalent(location: TalentLocation): TalentPicker<SpecType> {
		const talent = this.talents.find(talent => talent.getRow() == location.rowIdx && talent.getCol() == location.colIdx);
		if (!talent)
			throw new Error('No talent found with location: ' + location);
		return talent;
	}

	getTalentsString(): string {
		return this.talents.map(talent => String(talent.getPoints())).join('').replace(/0+$/g, '');
	}

	setTalentsString(str: string) {
		this.talents.forEach((talent, idx) => talent.setPoints(Number(str.charAt(idx)), false));
	}
}

class TalentPicker<SpecType extends Spec> extends Component {
	readonly config: TalentConfig<SpecType>;
	private readonly tree: TalentTreePicker<SpecType>;
	private readonly pointsDisplay: HTMLElement;

	private longTouchTimer?: number;

	constructor(parent: HTMLElement, player: Player<SpecType>, config: TalentConfig<SpecType>, tree: TalentTreePicker<SpecType>) {
		super(parent, 'talent-picker-root', document.createElement('a'));
		this.config = config;
		this.tree = tree;

		this.rootElem.style.gridRow = String(this.config.location.rowIdx + 1);
		this.rootElem.style.gridColumn = String(this.config.location.colIdx + 1);

		this.rootElem.dataset.maxPoints = String(this.config.maxPoints);
		this.rootElem.dataset.wowhead = 'noimage';

		this.pointsDisplay = document.createElement('span');
		this.pointsDisplay.classList.add('talent-picker-points');
		this.rootElem.appendChild(this.pointsDisplay);

		this.rootElem.addEventListener('click', event => {
			event.preventDefault();
		});
		this.rootElem.addEventListener('contextmenu', event => {
			event.preventDefault();
		});
		this.rootElem.addEventListener('touchstart', event => {
			event.preventDefault();
			this.longTouchTimer = setTimeout(() => {
				this.setPoints(0, true);
				this.tree.picker.update(TypedEvent.nextEventID());
				this.longTouchTimer = undefined;
			}, 750);
		});
		this.rootElem.addEventListener('touchend', event => {
			event.preventDefault();
			if (this.tree.picker.frozen)
				return;

			if (this.longTouchTimer != undefined) {
				clearTimeout(this.longTouchTimer);
				this.longTouchTimer = undefined;
			} else {
				return;
			}
			var newPoints = this.getPoints() + 1;
			if (this.config.maxPoints < newPoints) {
				newPoints = 0;
			}
			this.setPoints(newPoints, true);
			this.tree.picker.update(TypedEvent.nextEventID());
		});
		this.rootElem.addEventListener('mousedown', event => {
			if (this.tree.picker.frozen)
				return;

			const rightClick = isRightClick(event);
			if (rightClick) {
				this.setPoints(this.getPoints() - 1, true);
			} else {
				this.setPoints(this.getPoints() + 1, true);
			}
			this.tree.picker.update(TypedEvent.nextEventID());
		});
	}

	getRow(): number {
		return this.config.location.rowIdx;
	}

	getCol(): number {
		return this.config.location.colIdx;
	}

	getPoints(): number {
		const pts = Number(this.rootElem.dataset.points);
		return isNaN(pts) ? 0 : pts;
	}

	isFull(): boolean {
		return this.getPoints() >= this.config.maxPoints;
	}

	// Returns whether setting the points to newPoints would be a valid talent tree.
	canSetPoints(newPoints: number): boolean {
		const oldPoints = this.getPoints();

		if (newPoints > oldPoints) {
			const additionalPoints = newPoints - oldPoints;

			if (this.tree.picker.numPoints + additionalPoints > MAX_TALENT_POINTS) {
				return false;
			}

			if (this.tree.numPoints < this.getRow() * 5) {
				return false;
			}

			if (this.config.prereqLocation) {
				if (!this.tree.getTalent(this.config.prereqLocation).isFull())
					return false;
			}
		} else {
			const removedPoints = oldPoints - newPoints;

			// Figure out whether any lower talents would have the row requirement
			// broken by subtracting points.
			const pointTotalsByRow = [...Array(NUM_ROWS).keys()]
				.map(rowIdx => this.tree.talents.filter(talent => talent.getRow() == rowIdx))
				.map(talentsInRow => sum(talentsInRow.map(talent => talent.getPoints())));
			pointTotalsByRow[this.getRow()] -= removedPoints;

			const cumulativeTotalsByRow = pointTotalsByRow.map((_, rowIdx) => sum(pointTotalsByRow.slice(0, rowIdx + 1)));

			if (!this.tree.talents.every(talent =>
				talent.getPoints() == 0
				|| talent.getRow() == 0
				|| cumulativeTotalsByRow[talent.getRow() - 1] >= talent.getRow() * 5)) {
				return false;
			}

			if (this.config.prereqOfLocation) {
				if (this.tree.getTalent(this.config.prereqOfLocation).getPoints() > 0)
					return false;
			}
		}
		return true;
	}

	setPoints(newPoints: number, checkValidity: boolean) {
		const oldPoints = this.getPoints();
		newPoints = Math.max(0, newPoints);
		newPoints = Math.min(this.config.maxPoints, newPoints);

		if (checkValidity && !this.canSetPoints(newPoints))
			return;

		this.tree.numPoints += newPoints - oldPoints;
		this.rootElem.dataset.points = String(newPoints);

		this.pointsDisplay.textContent = newPoints + '/' + this.config.maxPoints;

		if (this.isFull()) {
			this.rootElem.classList.add('talent-full');
		} else {
			this.rootElem.classList.remove('talent-full');
		}

		const spellId = this.getSpellIdForPoints(newPoints);
		ActionId.fromSpellId(spellId).fill().then(actionId => {
			actionId.setWowheadHref(this.rootElem as HTMLAnchorElement);
			this.rootElem.style.backgroundImage = `url('${actionId.iconUrl}')`;
		});
	}

	getSpellIdForPoints(numPoints: number): number {
		// 0-indexed rank of talent
		const rank = Math.max(0, numPoints - 1);

		if (this.config.spellIds[rank]) {
			return this.config.spellIds[rank];
		} else {
			throw new Error('No rank ' + numPoints + ' for talent ' + this.config.fieldName);
		}
	}

	update() {
		if (this.canSetPoints(this.getPoints() + 1)) {
			this.rootElem.classList.add('talent-picker-can-add');
		} else {
			this.rootElem.classList.remove('talent-picker-can-add');
		}
	}
}

export type TalentsConfig<SpecType extends Spec> = Array<TalentTreeConfig<SpecType>>;

export type TalentTreeConfig<SpecType extends Spec> = {
	name: string;
	backgroundUrl: string;
	talents: Array<TalentConfig<SpecType>>;
};

export type TalentLocation = {
	// 0-indexed row in the tree
	rowIdx: number;
	// 0-indexed column in the tree
	colIdx: number;
};

export type TalentConfig<SpecType extends Spec> = {
	fieldName?: keyof SpecTalents<SpecType>

	location: TalentLocation;

	// Location of a prerequisite talent, if any
	prereqLocation?: TalentLocation;

	// Reverse of prereqLocation. This is populated automatically.
	prereqOfLocation?: TalentLocation;

	// Spell ID for each rank of this talent.
	// Omitted ranks will be inferred by incrementing from the last provided rank.
	spellIds: Array<number>;

	maxPoints: number;
};

export function newTalentsConfig<SpecType extends Spec>(talents: TalentsConfig<SpecType>): TalentsConfig<SpecType> {
	talents.forEach(tree => {
		tree.talents.forEach((talent, i) => {
			// Validate that talents are given in the correct order (left-to-right top-to-bottom).
			if (i != 0) {
				const prevTalent = tree.talents[i - 1];
				if (talent.location.rowIdx < prevTalent.location.rowIdx || (talent.location.rowIdx == prevTalent.location.rowIdx && talent.location.colIdx <= prevTalent.location.colIdx)) {
					throw new Error('Out-of-order talent: ' + talent.fieldName);
				}
			}

			// Infer omitted spell IDs.
			if (talent.spellIds.length < talent.maxPoints) {
				let curSpellId = talent.spellIds[talent.spellIds.length - 1];
				for (let pointIdx = talent.spellIds.length; pointIdx < talent.maxPoints; pointIdx++) {
					curSpellId++;
					talent.spellIds.push(curSpellId);
				}
			}
		});
	});
	return talents;
}
