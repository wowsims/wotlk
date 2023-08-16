import { Carousel, Tooltip } from 'bootstrap';
import { Component } from '../components/component.js';
import { Input, InputConfig } from '../components/input.js';
import { Class, Spec } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';
import { getSpecIcon } from '../proto_utils/utils.js';
import { TypedEvent } from '../typed_event.js';
import { isRightClick } from '../utils.js';
import { sum } from '../utils.js';
import { Player } from '../player.js';

const MAX_POINTS_PLAYER = 71;
const MAX_POINTS_HUNTER_PET = 16;
const MAX_POINTS_HUNTER_PET_BM = 20;

export interface TalentsPickerConfig<ModObject, TalentsProto> extends InputConfig<ModObject, string> {
	klass: Class,
	trees: TalentsConfig<TalentsProto>,
	pointsPerRow: number,
	maxPoints: number,
}

export class TalentsPicker<ModObject, TalentsProto> extends Input<ModObject, string> {
	private readonly config: TalentsPickerConfig<ModObject, TalentsProto>;

	readonly numRows: number;
	readonly numCols: number;
	readonly pointsPerRow: number;
	maxPoints: number;

	readonly trees: Array<TalentTreePicker<TalentsProto>>;

	constructor(parent: HTMLElement, modObject: ModObject, config: TalentsPickerConfig<ModObject, TalentsProto>) {
		super(parent, 'talents-picker-root', modObject, { ...config, inline: true });
		this.config = config;
		this.pointsPerRow = config.pointsPerRow;
		this.maxPoints = config.maxPoints;
		this.numRows = Math.max(...config.trees.map(treeConfig => treeConfig.talents.map(talentConfig => talentConfig.location.rowIdx).flat()).flat()) + 1;
		this.numCols = Math.max(...config.trees.map(treeConfig => treeConfig.talents.map(talentConfig => talentConfig.location.colIdx).flat()).flat()) + 1;
		this.rootElem.innerHTML = `
			<div id="talents-carousel" class="carousel slide">
				<div class="carousel-inner">
				</div>
				<button class="carousel-control-prev" type="button">
					<span class="carousel-control-prev-icon" aria-hidden="true"></span>
					<span class="visually-hidden">Previous</span>
				</button>
				<button class="carousel-control-next" type="button">
					<span class="carousel-control-next-icon" aria-hidden="true"></span>
					<span class="visually-hidden">Next</span>
				</button>
			</div>
		`
		const carouselContainer = this.rootElem.querySelector('.carousel-inner') as HTMLElement;
		const carouselPrevBtn = this.rootElem.querySelector('.carousel-control-prev') as HTMLButtonElement;
		const carouselNextBtn = this.rootElem.querySelector('.carousel-control-next') as HTMLButtonElement;

		this.trees = config.trees.map((treeConfig, i) => {
			const carouselItem = document.createElement('div');
			carouselContainer.appendChild(carouselItem);

			carouselItem.classList.add('carousel-item');
			// Set middle talents active by default for mobile slider
			if (i === 1) carouselItem.classList.add('active');

			// If using a hunter pet, add 3 to skip the hunter specs
			if (treeConfig.name === 'Ferocity') i += 3;
			if (treeConfig.name === 'Tenacity') i += 4;
			if (treeConfig.name === 'Cunning') i += 5;

			return new TalentTreePicker(carouselItem, treeConfig, this, config.klass, i);
		});
		this.trees.forEach(tree => tree.talents.forEach(talent => talent.setPoints(0, false)));

		let carouselitemIdx = 0;
		const slidePrev = () => {
			if (carouselitemIdx >= 1) return;
			carouselitemIdx += 1;
			carouselContainer.style.transform = `translateX(${33.3 * carouselitemIdx}%)`;
			carouselContainer.children[Math.abs(carouselitemIdx - 2) % 3]!.classList.remove('active');
			carouselContainer.children[Math.abs(carouselitemIdx - 1) % 3]!.classList.add('active')
		}
		const slideNext = () => {
			if (carouselitemIdx <= -1) return;
			carouselitemIdx -= 1;
			carouselContainer.style.transform = `translateX(${33.3 * carouselitemIdx}%)`;
			carouselContainer.children[Math.abs(carouselitemIdx) % 3]!.classList.remove('active');
			carouselContainer.children[Math.abs(carouselitemIdx) + 1 % 3]!.classList.add('active');
		}

		carouselPrevBtn.addEventListener('click', slidePrev);
		carouselPrevBtn.addEventListener('touchend', slidePrev);
		carouselNextBtn.addEventListener('click', slideNext);
		carouselNextBtn.addEventListener('touchend', slideNext);

		this.init();
	}

	getInputElem(): HTMLElement {
		return this.rootElem;
	}

	getInputValue(): string {
		return this.trees.map(tree => tree.getTalentsString()).join('-').replace(/-+$/g, '');
	}

	setInputValue(newValue: string) {
		const parts = newValue.split('-');
		this.trees.forEach((tree, idx) => tree.setTalentsString(parts[idx] || ''));
		this.updateTrees();
	}

	updateTrees() {
		if (this.isFull()) {
			this.rootElem.classList.add('talents-full');
		} else {
			this.rootElem.classList.remove('talents-full');
		}
		this.trees.forEach(tree => tree.update());
	}

	get numPoints() {
		return sum(this.trees.map(tree => tree.numPoints));
	}

	isFull() {
		return this.numPoints >= this.maxPoints;
	}

	setMaxPoints(newMaxPoints: number) {
		if (newMaxPoints != this.maxPoints) {
			this.maxPoints = newMaxPoints;
			this.updateTrees();
		}
	}

	isHunterPet(): boolean  {
		return ['Cunning', 'Ferocity', 'Tenacity'].includes(this.config.trees[0].name)
	}
}

class TalentTreePicker<TalentsProto> extends Component {
	private readonly config: TalentTreeConfig<TalentsProto>;
	private readonly title: HTMLElement;
	private readonly pointsElem: HTMLElement;

	readonly talents: Array<TalentPicker<TalentsProto>>;
	readonly picker: TalentsPicker<any, TalentsProto>;

	// The current number of points in this tree
	numPoints: number;

	constructor(parent: HTMLElement, config: TalentTreeConfig<TalentsProto>, picker: TalentsPicker<any, TalentsProto>, klass: Class, specNumber: number) {
		super(parent, 'talent-tree-picker-root');
		this.config = config;
		this.numPoints = 0;
		this.picker = picker;

		this.rootElem.innerHTML = `
			<div class="talent-tree-header">
				<img src=${getSpecIcon(klass, specNumber)} class="talent-tree-icon" />
				<span class="talent-tree-title"></span>
				<span class="talent-tree-points"></span>
				<button
					class="talent-tree-reset btn btn-link link-danger"
					data-bs-title="Reset talent points"
					data-bs-toggle="tooltip"
				>
					<i class="fa fa-times"></i>
				</button>
			</div>
			<div class="talent-tree-background"></div>
			<div class="talent-tree-main"></div>
    `;

		this.title = this.rootElem.getElementsByClassName('talent-tree-title')[0] as HTMLElement;
		this.pointsElem = this.rootElem.querySelector('.talent-tree-points') as HTMLElement;

		const background = this.rootElem.querySelector('.talent-tree-background') as HTMLElement;
		background.style.backgroundImage = `url('${config.backgroundUrl}')`;

		const main = this.rootElem.querySelector('.talent-tree-main') as HTMLElement;
		main.style.gridTemplateRows = `repeat(${this.picker.numRows}, 1fr)`;
		// Add 2 for spacing on the sides
		main.style.gridTemplateColumns = `repeat(${this.picker.numCols}, 1fr)`;

		const iconSize = '3.5rem'
		main.style.height = `calc(${iconSize} * ${this.picker.numRows})`
		main.style.maxWidth = `calc(${iconSize} * ${this.picker.numCols})`
		this.rootElem.style.maxWidth = `calc(${iconSize} * ${this.picker.numCols + 2})`

		this.talents = config.talents.map(talent => new TalentPicker(main, talent, this));
		this.talents.forEach(talent => {
			if (talent.config.prereqLocation) {
				this.getTalent(talent.config.prereqLocation).config.prereqOfLocation = talent.config.location;
			}
		});

		const resetBtn = this.rootElem.querySelector('.talent-tree-reset') as HTMLElement;
		new Tooltip(resetBtn)
		resetBtn.addEventListener('click', event => {
			this.talents.forEach(talent => talent.setPoints(0, false));
			this.picker.inputChanged(TypedEvent.nextEventID());
		});
	}

	update() {
		this.title.innerHTML = this.config.name
		this.pointsElem.textContent = `${this.numPoints} / ${this.getMaxSpendablePoints()}`
		this.talents.forEach(talent => talent.update());
	}

	getTalent(location: TalentLocation): TalentPicker<TalentsProto> {
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

	getMaxSpendablePoints() {
		if (!this.picker.isHunterPet()) return MAX_POINTS_PLAYER;
		if ((this.picker.modObject as Player<Spec.SpecHunter>).getTalents().beastMastery) return MAX_POINTS_HUNTER_PET_BM;
		return MAX_POINTS_HUNTER_PET;
	}
}

class TalentPicker<TalentsProto> extends Component {
	readonly config: TalentConfig<TalentsProto>;
	private readonly tree: TalentTreePicker<TalentsProto>;
	private readonly pointsDisplay: HTMLElement;

	private longTouchTimer?: number;

	constructor(parent: HTMLElement, config: TalentConfig<TalentsProto>, tree: TalentTreePicker<TalentsProto>) {
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
			this.longTouchTimer = window.setTimeout(() => {
				this.setPoints(0, true);
				this.tree.picker.inputChanged(TypedEvent.nextEventID());
				this.longTouchTimer = undefined;
			}, 750);
		});
		this.rootElem.addEventListener('touchend', event => {
			event.preventDefault();
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
			this.tree.picker.inputChanged(TypedEvent.nextEventID());
		});
		this.rootElem.addEventListener('mousedown', event => {
			const rightClick = isRightClick(event);
			if (rightClick) {
				this.setPoints(this.getPoints() - 1, true);
			} else {
				this.setPoints(this.getPoints() + 1, true);
			}
			this.tree.picker.inputChanged(TypedEvent.nextEventID());
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

			if (this.tree.picker.numPoints + additionalPoints > this.tree.picker.maxPoints) {
				return false;
			}

			if (this.tree.numPoints < this.getRow() * this.tree.picker.pointsPerRow) {
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
			const pointTotalsByRow = [...Array(this.tree.picker.numRows).keys()]
				.map(rowIdx => this.tree.talents.filter(talent => talent.getRow() == rowIdx))
				.map(talentsInRow => sum(talentsInRow.map(talent => talent.getPoints())));
			pointTotalsByRow[this.getRow()] -= removedPoints;

			const cumulativeTotalsByRow = pointTotalsByRow.map((_, rowIdx) => sum(pointTotalsByRow.slice(0, rowIdx + 1)));

			if (!this.tree.talents.every(talent =>
				talent.getPoints() == 0
				|| talent.getRow() == 0
				|| cumulativeTotalsByRow[talent.getRow() - 1] >= talent.getRow() * this.tree.picker.pointsPerRow)) {
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
			throw new Error(`No rank ${numPoints} for talent ${String(this.config.fieldName)}`);
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

export type TalentsConfig<TalentsProto> = Array<TalentTreeConfig<TalentsProto>>;

export type TalentTreeConfig<TalentsProto> = {
	name: string;
	backgroundUrl: string;
	talents: Array<TalentConfig<TalentsProto>>;
};

export type TalentLocation = {
	// 0-indexed row in the tree
	rowIdx: number;
	// 0-indexed column in the tree
	colIdx: number;
};

export type TalentConfig<TalentsProto> = {
	fieldName?: keyof TalentsProto | string;

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

export function newTalentsConfig<TalentsProto>(talents: TalentsConfig<TalentsProto>): TalentsConfig<TalentsProto> {
	talents.forEach(tree => {
		tree.talents.forEach((talent, i) => {
			// Validate that talents are given in the correct order (left-to-right top-to-bottom).
			if (i != 0) {
				const prevTalent = tree.talents[i - 1];
				if (talent.location.rowIdx < prevTalent.location.rowIdx || (talent.location.rowIdx == prevTalent.location.rowIdx && talent.location.colIdx <= prevTalent.location.colIdx)) {
					throw new Error(`Out-of-order talent: ${String(talent.fieldName)}`);
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
