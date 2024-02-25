// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element, fragment, ref } from 'tsx-vanilla';
import { Tooltip } from 'bootstrap';

import { Component } from '../components/component.js';
import { CopyButton } from '../components/copy_button.js';
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

export interface TalentsPickerConfig<TalentsProto> extends InputConfig<Player<Spec>, string> {
	klass: Class,
	maxPoints: number,
	pointsPerRow: number,
	trees: TalentsConfig<TalentsProto>,
}

export class TalentsPicker<TalentsProto> extends Input<Player<Spec>, string> {
	private readonly config: TalentsPickerConfig<TalentsProto>;

	readonly numRows: number;
	readonly numCols: number;
	readonly pointsPerRow: number;
	maxPoints: number;

	readonly trees: Array<TalentTreePicker<TalentsProto>>;

	constructor(parent: HTMLElement, player: Player<Spec>, config: TalentsPickerConfig<TalentsProto>) {
		super(parent, 'talents-picker-root', player, { ...config, inline: true });
		this.config = config;
		this.pointsPerRow = config.pointsPerRow;
		this.numRows = Math.max(...config.trees.map(treeConfig => treeConfig.talents.map(talentConfig => talentConfig.location.rowIdx).flat()).flat()) + 1;
		this.numCols = Math.max(...config.trees.map(treeConfig => treeConfig.talents.map(talentConfig => talentConfig.location.colIdx).flat()).flat()) + 1;
		this.maxPoints = config.maxPoints

		const pointsRemainingElemRef = ref<HTMLSpanElement>();
		const getPointsRemaining = () => this.maxPoints - player.getTalentTreePoints().reduce((sum, points) => sum + points, 0);

		const PointsRemainingElem = () => {
			const pointsRemaining = getPointsRemaining();
			return <span className="talent-tree-points" ref={pointsRemainingElemRef}>{pointsRemaining}</span>
		}

		TypedEvent.onAny([player.talentsChangeEmitter]).on(() => {
			pointsRemainingElemRef.value!.replaceWith(PointsRemainingElem())
		});

		const actionsContainerRef = ref<HTMLDivElement>();
		this.rootElem.appendChild(
			<div id="talents-carousel" className="carousel slide">
				<div className="talents-picker-header">
					<div>
						<label>Points Remaining:</label>
						{PointsRemainingElem()}
					</div>
					<div className="talents-picker-actions" ref={actionsContainerRef}></div>
				</div>
				<div className="carousel-inner">
				</div>
				<div id="talents-carousel" className="carousel slide">
					<div className="carousel-inner">
					</div>
					<button className="carousel-control-prev" type="button">
						<span className="carousel-control-prev-icon" attributes={{'aria-hidden':true}}></span>
						<span className="visually-hidden">Previous</span>
					</button>
					<button className="carousel-control-next" type="button">
						<span className="carousel-control-next-icon" attributes={{'aria-hidden':true}}></span>
						<span className="visually-hidden">Next</span>
					</button>
				</div>
			</div>
		);

		new CopyButton(actionsContainerRef.value!, {
			extraCssClasses: ['btn-sm', 'btn-outline-primary', 'copy-talents'],
			getContent: () => player.getTalentsString(),
			text: "Copy",
			tooltip: "Copy talent string",
		});

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
		carouselNextBtn.addEventListener('click', slideNext);

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
	readonly picker: TalentsPicker<TalentsProto>;

	// The current number of points in this tree
	numPoints: number;

	constructor(parent: HTMLElement, config: TalentTreeConfig<TalentsProto>, picker: TalentsPicker<TalentsProto>, klass: Class, specNumber: number) {
		super(parent, 'talent-tree-picker-root');
		this.config = config;
		this.numPoints = 0;
		this.picker = picker;

		this.rootElem.appendChild(
			<>
				<div className="talent-tree-header">
					<img src={getSpecIcon(klass, specNumber)} className="talent-tree-icon" />
					<span className="talent-tree-title"></span>
					<span className="talent-tree-points"></span>
					<button className="talent-tree-reset btn btn-link link-danger">
						<i className="fa fa-times"></i>
					</button>
				</div>
				<div className="talent-tree-background"></div>
				<div className="talent-tree-main"></div>
			</>
		);

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
		// Process parent<->child mapping
		this.talents.forEach(talent => {
			if (talent.config.prereqLocation) {
				this.getTalent(talent.config.prereqLocation).config.childLocations!.push(talent.config.location);
			}
		});
		// Loop through all and have talent add in divs/items for child dependencies
		// It'd be nicer to have this in talent constructor but json would have to be updated
		const recurseCalcIdx = (t: TalentPicker<TalentsProto>, z: number) => {
			t.initChildReqs();
			t.zIndex = z;
			for (const cl of t.config.childLocations!) {
				const c = this.getTalent(cl);
				c.parentReq = t.getChildReqArrow(cl);
				recurseCalcIdx(c, z-2);
			}
		}
		// Start at top of each heirachy chain and recurse down
		for (const t of this.talents) {
			if (t.config.childLocations!.length == 0)
				continue;
			if (t.config.prereqLocation !== undefined)
				continue;
			recurseCalcIdx(t, 20);
		}
		const resetBtn = this.rootElem.querySelector('.talent-tree-reset') as HTMLElement;
		new Tooltip(resetBtn, {
			title: 'Reset talent points',
		});
		resetBtn.addEventListener('click', _event => {
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


type ReqDir = 'down' | 'right' | 'left' | 'rightdown' | 'leftdown';
class TalentReqArrow extends Component {
	private dir: ReqDir;
	private zIdx: number;
	readonly parentLoc: TalentLocation;
	readonly childLoc: TalentLocation;

	constructor(parent: HTMLElement, parentLoc: TalentLocation, childLoc: TalentLocation) {
		super(parent, 'talent-picker-req-arrow', document.createElement('div'));
		this.zIdx = 0;
		this.parentLoc = parentLoc;
		this.childLoc = childLoc;

		this.rootElem.style.gridRow = String(parentLoc.rowIdx + 1);
		this.rootElem.style.gridColumn = String(parentLoc.colIdx + 1);

		let rowEnd = Math.max(parentLoc.rowIdx, childLoc.rowIdx) + 1;
		let colEnd = Math.max(parentLoc.colIdx, childLoc.colIdx) + 1;

		// Calculate where we need to 'point'
		if (parentLoc.rowIdx == childLoc.rowIdx) {
			this.dir = parentLoc.colIdx < childLoc.colIdx ? 'right' : 'left';
			this.rootElem.dataset.reqArrowColSize = String(Math.abs(parentLoc.colIdx - childLoc.colIdx));
			colEnd = this.dir == 'left' ? colEnd+1 : colEnd-1;
		} else {
			if (parentLoc.colIdx == childLoc.colIdx) {
				this.dir = 'down';
				this.rootElem.dataset.reqArrowRowSize = String(Math.abs(parentLoc.rowIdx - childLoc.rowIdx));
				rowEnd += 1;
			}
			else {
				this.dir = parentLoc.colIdx < childLoc.colIdx ? 'rightdown' : 'leftdown';
				this.rootElem.dataset.reqArrowColSize = String(Math.abs(parentLoc.colIdx - childLoc.colIdx));
				this.rootElem.dataset.reqArrowRowSize = String(Math.abs(parentLoc.rowIdx - childLoc.rowIdx));
				rowEnd += 1;
				colEnd = this.dir == 'rightdown' ? colEnd+1 : colEnd-1;
				this.rootElem.appendChild(<div></div>)
			}
		}

		this.rootElem.style.gridRowEnd = String(rowEnd);
		this.rootElem.style.gridColumnEnd = String(colEnd);
		this.rootElem.classList.add(`talent-picker-req-arrow-${this.dir}`);
	}

	get zIndex() {
		return this.zIdx;
	}

	set zIndex(z: number) {
		this.zIdx = z;
		this.rootElem.style.zIndex = String(z);
	}

	setReqFufilled(isFufilled: boolean) {
		if (isFufilled)
			this.rootElem.dataset.reqActive = 'true';
		else
			delete this.rootElem.dataset.reqActive;
	}
}

class TalentPicker<TalentsProto> extends Component {
	readonly config: TalentConfig<TalentsProto>;
	private readonly tree: TalentTreePicker<TalentsProto>;
	private readonly pointsDisplay: HTMLElement;

	private longTouchTimer?: number;
	private childReqs: TalentReqArrow[];
	private zIdx: number;
	parentReq: TalentReqArrow | null;

	constructor(parent: HTMLElement, config: TalentConfig<TalentsProto>, tree: TalentTreePicker<TalentsProto>) {
		super(parent, 'talent-picker-root', document.createElement('a'));
		this.config = config;
		this.tree = tree;
		this.childReqs = [];
		this.parentReq = null;
		this.zIdx = 0;

		this.rootElem.style.gridRow = String(this.config.location.rowIdx + 1);
		this.rootElem.style.gridColumn = String(this.config.location.colIdx + 1);

		this.rootElem.dataset.maxPoints = String(this.config.maxPoints);
		this.rootElem.dataset.whtticon = 'false';

		this.pointsDisplay = document.createElement('span');
		this.pointsDisplay.classList.add('talent-picker-points');
		this.rootElem.appendChild(this.pointsDisplay);

		this.rootElem.addEventListener('click', event => {
			event.preventDefault();
		});
		this.rootElem.addEventListener('contextmenu', event => {
			event.preventDefault();
		});
		this.rootElem.addEventListener('touchmove', _event => {
			if (this.longTouchTimer != undefined) {
				clearTimeout(this.longTouchTimer);
				this.longTouchTimer = undefined;
			}
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

	initChildReqs(): void {
		if (this.config.childLocations!.length == 0)
			return;

		for (const c of this.config.childLocations!) {
			this.childReqs.push(new TalentReqArrow(this.rootElem.parentElement!, this.config.location, c));
		}
	}

	getChildReqArrow(loc: TalentLocation): TalentReqArrow {
		for (let c of this.childReqs) {
			if (c.childLoc === loc) {
				return c;
			}
		}
		throw Error("missing child prereq?");
	}

	get zIndex() {
		return this.zIdx;
	}

	set zIndex(z: number) {
		this.zIdx = z;
		this.rootElem.style.zIndex = String(this.zIdx);

		for (const c of this.childReqs) {
			c.zIndex = this.zIdx-1;
		}
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

			for (const c of this.config.childLocations!) {
				if (this.tree.getTalent(c).getPoints() > 0)
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
		let canSetPoints = this.canSetPoints(this.getPoints() + 1);
		if (canSetPoints) {
			this.rootElem.classList.add('talent-picker-can-add');
		} else {
			this.rootElem.classList.remove('talent-picker-can-add');
		}

		if (this.parentReq) {
			this.parentReq.setReqFufilled(canSetPoints || this.isFull());
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

	// Child talents depending on this talent. This is populated automatically.
	childLocations?: TalentLocation[];

	// Spell ID for each rank of this talent.
	// Omitted ranks will be inferred by incrementing from the last provided rank.
	spellIds: Array<number>;

	maxPoints: number;
};

export function newTalentsConfig<TalentsProto>(talents: TalentsConfig<TalentsProto>): TalentsConfig<TalentsProto> {
	talents.forEach(tree => {
		tree.talents.forEach((talent, i) => {
			talent.childLocations = [];
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
