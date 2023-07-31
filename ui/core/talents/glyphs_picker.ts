import { Glyphs } from '../proto/common.js';
import { ItemQuality } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';
import { setItemQualityCssClass } from '../css_utils.js';
import { Player } from '../player.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { stringComparator } from '../utils.js';

import { Component } from '../components/component.js';
import { Input } from '../components/input.js';
import { BaseModal } from '../components/base_modal.js';

export type GlyphConfig = {
	name: string,
	description: string,
	iconUrl: string,
};

export type GlyphsConfig = {
	majorGlyphs: Record<number, GlyphConfig>,
	minorGlyphs: Record<number, GlyphConfig>,
};

interface GlyphData {
	id: number,
	name: string,
	description: string,
	iconUrl: string,
	quality: ItemQuality | null,
}

const emptyGlyphData: GlyphData = {
	id: 0,
	name: 'Empty',
	description: '',
	iconUrl: 'https://wow.zamimg.com/images/wow/icons/medium/inventoryslot_empty.jpg',
	quality: null,
};

export class GlyphsPicker extends Component {
	private readonly glyphsConfig: GlyphsConfig;

	majorGlyphPickers: Array<GlyphPicker> = [];
	minorGlyphPickers: Array<GlyphPicker> = [];

	constructor(parent: HTMLElement, player: Player<any>, glyphsConfig: GlyphsConfig) {
		super(parent, 'glyphs-picker-root');
		this.glyphsConfig = glyphsConfig;

		const majorGlyphs = Object.keys(glyphsConfig.majorGlyphs).map(idStr => Number(idStr));
		const minorGlyphs = Object.keys(glyphsConfig.minorGlyphs).map(idStr => Number(idStr));

		const majorGlyphsData = majorGlyphs.map(glyph => this.getGlyphData(glyph));
		const minorGlyphsData = minorGlyphs.map(glyph => this.getGlyphData(glyph));

		majorGlyphsData.sort((a, b) => stringComparator(a.name, b.name));
		minorGlyphsData.sort((a, b) => stringComparator(a.name, b.name));

		this.majorGlyphPickers = (['major1', 'major2', 'major3'] as Array<keyof Glyphs>).map(glyphField => new GlyphPicker(this.rootElem, player, majorGlyphsData, glyphField, true));
		this.minorGlyphPickers = (['minor1', 'minor2', 'minor3'] as Array<keyof Glyphs>).map(glyphField => new GlyphPicker(this.rootElem, player, minorGlyphsData, glyphField, false));
	}

	// In case we ever want to parse description from tooltip HTML.
	//static descriptionRegex = /<a href=\\"\/wotlk.*>(.*)<\/a>/g;
	getGlyphData(glyph: number): GlyphData {
		const glyphConfig = this.glyphsConfig.majorGlyphs[glyph] || this.glyphsConfig.minorGlyphs[glyph];

		return {
			id: glyph,
			name: glyphConfig.name,
			description: glyphConfig.description,
			iconUrl: glyphConfig.iconUrl,
			quality: ItemQuality.ItemQualityCommon,
		};
	}
}

class GlyphPicker extends Input<Player<any>, number> {
	readonly player: Player<any>;
	private readonly iconElem: HTMLAnchorElement;

	private readonly glyphOptions: Array<GlyphData>;
	selectedGlyph: GlyphData;

	constructor(parent: HTMLElement, player: Player<any>, glyphOptions: Array<GlyphData>, glyphField: keyof Glyphs, isMajor: boolean) {
		super(parent, 'glyph-picker-root', player, {
			changedEvent: (player: Player<any>) => player.glyphsChangeEmitter,
			getValue: (player: Player<any>) => player.getGlyphs()[glyphField] as number,
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const glyphs = player.getGlyphs();
				(glyphs[glyphField] as number) = newValue;
				player.setGlyphs(eventID, glyphs);
			},
		});
		if (!isMajor) {
			this.rootElem.classList.add('minor');
		}
		this.player = player;
		this.glyphOptions = glyphOptions;
		this.selectedGlyph = emptyGlyphData;

		this.rootElem.innerHTML = `<a class="glyph-picker-icon"></a>`;

		this.iconElem = this.rootElem.getElementsByClassName('glyph-picker-icon')[0] as HTMLAnchorElement;
		this.iconElem.addEventListener('click', event => {
			event.preventDefault();
			const selectorModal = new GlyphSelectorModal(this.rootElem.closest('.individual-sim-ui')!, this, this.glyphOptions);
		});
		this.iconElem.addEventListener('touchstart', event => {
			event.preventDefault();
			const selectorModal = new GlyphSelectorModal(this.rootElem.closest('.individual-sim-ui')!, this, this.glyphOptions);
		});
		this.iconElem.addEventListener('touchend', event => {
			event.preventDefault();
		});

		this.init();
	}

	getInputElem(): HTMLElement {
		return this.iconElem;
	}

	getInputValue(): number {
		return this.selectedGlyph.id;
	}

	setInputValue(newValue: number) {
		this.selectedGlyph = this.glyphOptions.find(glyphData => glyphData.id == newValue) || emptyGlyphData;

		this.iconElem.style.backgroundImage = `url('${this.selectedGlyph.iconUrl}')`;
		this.iconElem.href = this.selectedGlyph.id == 0 ? '' : ActionId.makeItemUrl(this.selectedGlyph.id);
	}
}

class GlyphSelectorModal extends BaseModal {
	constructor(parent: HTMLElement, glyphPicker: GlyphPicker, glyphOptions: Array<GlyphData>) {
		super(parent, 'glyph-modal', { title: 'Glyphs' });

		this.body.innerHTML = `
			<div class="input-root">
				<input class="selector-modal-search form-control" type="text" placeholder="Search...">
			</div>
			<ul class="selector-modal-list"></ul>
		`;

		const listElem = this.rootElem.getElementsByClassName('selector-modal-list')[0] as HTMLElement;

		glyphOptions = [emptyGlyphData].concat(glyphOptions);
		const listItemElems = glyphOptions.map((glyphData, glyphIdx) => {
			const listItemElem = document.createElement('li');
			listItemElem.classList.add('selector-modal-list-item');
			listElem.appendChild(listItemElem);

			listItemElem.dataset.idx = String(glyphIdx);

			listItemElem.innerHTML = `
        <a class="selector-modal-list-item-icon"></a>
        <a class="selector-modal-list-item-name">${glyphData.name}</a>
				<span class="selector-modal-list-item-description">${glyphData.description}</span>
      `;

			(listItemElem.children[0] as HTMLAnchorElement).href = glyphData.id == 0 ? '' : ActionId.makeItemUrl(glyphData.id);
			(listItemElem.children[1] as HTMLAnchorElement).href = glyphData.id == 0 ? '' : ActionId.makeItemUrl(glyphData.id);
			const iconElem = listItemElem.getElementsByClassName('selector-modal-list-item-icon')[0] as HTMLImageElement;
			iconElem.style.backgroundImage = `url('${glyphData.iconUrl}')`;

			const nameElem = listItemElem.getElementsByClassName('selector-modal-list-item-name')[0] as HTMLImageElement;
			setItemQualityCssClass(nameElem, glyphData.quality);

			const onclick = (event: Event) => {
				event.preventDefault();
				glyphPicker.setValue(TypedEvent.nextEventID(), glyphData.id);
			};
			nameElem.addEventListener('click', onclick);
			iconElem.addEventListener('click', onclick);

			return listItemElem;
		});

		const updateSelected = () => {
			const selectedGlyphId = glyphPicker.selectedGlyph.id;

			listItemElems.forEach(elem => {
				const listItemIdx = parseInt(elem.dataset.idx!);
				const listItemData = glyphOptions[listItemIdx];

				elem.classList.remove('active');
				if (listItemData.id == selectedGlyphId) {
					elem.classList.add('active');
				}
			});
		};
		updateSelected();

		const applyFilters = () => {
			let validItemElems = listItemElems;
			const selectedGlyph = glyphPicker.selectedGlyph;

			validItemElems = validItemElems.filter(elem => {
				const listItemIdx = parseInt(elem.dataset.idx!);
				const listItemData = glyphOptions[listItemIdx];

				if (searchInput.value.length > 0) {
					const searchQuery = searchInput.value.toLowerCase().split(" ");
					const name = listItemData.name.toLowerCase();

					var include = true;
					searchQuery.forEach(v => {
						if (!name.includes(v))
							include = false;
					});
					if (!include) {
						return false;
					}
				}

				return true;
			});

			let numShown = 0;
			listItemElems.forEach(elem => {
				if (validItemElems.includes(elem)) {
					elem.classList.remove('hidden');
					numShown++;
					if (numShown % 2 == 0) {
						elem.classList.remove('odd');
					} else {
						elem.classList.add('odd');
					}
				} else {
					elem.classList.add('hidden');
				}
			});
		};

		const searchInput = this.rootElem.getElementsByClassName('selector-modal-search')[0] as HTMLInputElement;
		searchInput.addEventListener('input', applyFilters);
		searchInput.addEventListener("keyup", ev => {
			if (ev.key == "Enter") {
				listItemElems.find(ele => {
					if (ele.classList.contains("hidden")) {
						return false;
					}
					const nameElem = ele.getElementsByClassName('selector-modal-list-item-name')[0] as HTMLElement;
					nameElem.click();
					return true;
				});
			}
		});

		glyphPicker.player.glyphsChangeEmitter.on(() => {
			applyFilters();
			updateSelected();
		});
	}
}
