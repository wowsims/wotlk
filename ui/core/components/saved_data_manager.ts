import { EventID, TypedEvent } from '../typed_event.js';
import { ContentBlock, ContentBlockHeaderConfig } from './content_block';

import { Component } from '../components/component.js';
import { Tooltip } from 'bootstrap';

export type SavedDataManagerConfig<ModObject, T> = {
	label: string;
	header?: ContentBlockHeaderConfig;
	presetsOnly?: boolean;
	storageKey: string;
	changeEmitters: Array<TypedEvent<any>>,
	equals: (a: T, b: T) => boolean;
	getData: (modObject: ModObject) => T;
	setData: (eventID: EventID, modObject: ModObject, data: T) => void;
	toJson: (a: T) => any;
	fromJson: (obj: any) => T;
};

export type SavedDataConfig<ModObject, T> = {
	name: string;
	data: T;
	tooltip?: string;
	isPreset?: boolean;

	// If set, will automatically hide the saved data when this evaluates to false.
	enableWhen?: (obj: ModObject) => boolean;
};

type SavedData<ModObject, T> = {
	name: string;
	data: T;
	elem: HTMLElement;
	enableWhen?: (obj: ModObject) => boolean;
};

export class SavedDataManager<ModObject, T> extends Component {
	private readonly modObject: ModObject;
	private readonly config: SavedDataManagerConfig<ModObject, T>;

	private readonly userData: Array<SavedData<ModObject, T>>;
	private readonly presets: Array<SavedData<ModObject, T>>;

	private readonly savedDataDiv: HTMLElement
	private readonly presetDataDiv: HTMLElement;
	private readonly customDataDiv: HTMLElement;
	private readonly saveInput?: HTMLInputElement;

	private frozen: boolean;

	constructor(parent: HTMLElement, modObject: ModObject, config: SavedDataManagerConfig<ModObject, T>) {
		super(parent, 'saved-data-manager-root');
		this.modObject = modObject;
		this.config = config;

		this.userData = [];
		this.presets = [];
		this.frozen = false;

		let contentBlock = new ContentBlock(this.rootElem, 'saved-data', { header: config.header });

		contentBlock.bodyElement.innerHTML = `
			<div class="saved-data-container hide">
				<div class="saved-data-presets"></div>
				<div class="saved-data-custom"></div>
			</div>
		`;
		this.savedDataDiv = contentBlock.bodyElement.querySelector('.saved-data-container') as HTMLElement;
		this.presetDataDiv = contentBlock.bodyElement.querySelector('.saved-data-presets') as HTMLElement;
		this.customDataDiv = contentBlock.bodyElement.querySelector('.saved-data-custom') as HTMLElement;

		if (!config.presetsOnly) {
			contentBlock.bodyElement.appendChild(this.buildCreateContainer());
			this.saveInput = contentBlock.bodyElement.querySelector('.saved-data-save-input') as HTMLInputElement;
		}
	}

	addSavedData(config: SavedDataConfig<ModObject, T>) {
		this.savedDataDiv.classList.remove('hide');

		const newData = this.makeSavedData(config);
		const dataArr = config.isPreset ? this.presets : this.userData;
		const oldIdx = dataArr.findIndex(data => data.name == config.name);

		if (oldIdx == -1) {
			if (config.isPreset) {
				this.presetDataDiv.appendChild(newData.elem);
			} else {
				this.customDataDiv.appendChild(newData.elem);
			}
			dataArr.push(newData);
		} else {
			dataArr[oldIdx].elem.replaceWith(newData.elem)
			dataArr[oldIdx] = newData;
		}
	}

	private makeSavedData(config: SavedDataConfig<ModObject, T>): SavedData<ModObject, T> {
		const dataElemFragment = document.createElement('fragment');
		dataElemFragment.innerHTML = `
			<div class="saved-data-set-chip badge rounded-pill">
				<a href="javascript:void(0)" class="saved-data-set-name" role="button">${config.name}</a>
			</div>
		`;

		const dataElem = dataElemFragment.children[0] as HTMLElement;
		dataElem.addEventListener('click', event => {
			this.config.setData(TypedEvent.nextEventID(), this.modObject, config.data);

			if (this.saveInput)
				this.saveInput.value = config.name;
		});

		if (!config.isPreset) {
			let deleteFragment = document.createElement('fragment');
			deleteFragment.innerHTML = `
				<a
					href="javascript:void(0)"
					class="saved-data-set-delete"
					role="button"
				>
					<i class="fa fa-times fa-lg"></i>
				</a>
			`;

			const deleteButton = deleteFragment.children[0] as HTMLElement;
			dataElem.appendChild(deleteButton);

			const tooltip = Tooltip.getOrCreateInstance(deleteButton, {title:`Delete saved ${this.config.label}`});

			deleteButton.addEventListener('click', event => {
				event.stopPropagation();
				const shouldDelete = confirm(`Delete saved ${this.config.label} '${config.name}'?`);
				if (!shouldDelete)
					return;

				tooltip.dispose();

				const idx = this.userData.findIndex(data => data.name == config.name);
				this.userData[idx].elem.remove();
				this.userData.splice(idx, 1);
				this.saveUserData();
			});
		}

		if (config.tooltip) {
			Tooltip.getOrCreateInstance(dataElem, {
				title: config.tooltip,
				placement: 'bottom',
				html: true,
			});
		}

		const checkActive = () => {
			if (this.config.equals(config.data, this.config.getData(this.modObject))) {
				dataElem.classList.add('active');
			} else {
				dataElem.classList.remove('active');
			}

			if (config.enableWhen && !config.enableWhen(this.modObject)) {
				dataElem.classList.add('disabled');
			} else {
				dataElem.classList.remove('disabled');
			}
		};

		checkActive();
		this.config.changeEmitters.forEach(emitter => emitter.on(checkActive));

		return {
			name: config.name,
			data: config.data,
			elem: dataElem,
			enableWhen: config.enableWhen,
		};
	}

	// Save data to window.localStorage.
	private saveUserData() {
		const userData: Record<string, Object> = {};
		this.userData.forEach(savedData => {
			userData[savedData.name] = this.config.toJson(savedData.data);
		});

		if (this.userData.length == 0 && this.presets.length == 0)
			this.savedDataDiv.classList.add('hide');

		window.localStorage.setItem(this.config.storageKey, JSON.stringify(userData));
	}

	// Load data from window.localStorage.
	loadUserData() {
		const dataStr = window.localStorage.getItem(this.config.storageKey);
		if (!dataStr)
			return;

		let jsonData;
		try {
			jsonData = JSON.parse(dataStr);
		} catch (e) {
			console.warn('Invalid json for local storage value: ' + dataStr);
		}

		for (let name in jsonData) {
			try {
				this.addSavedData({
					name: name,
					data: this.config.fromJson(jsonData[name]),
				});
			} catch (e) {
				console.warn('Failed parsing saved data: ' + jsonData[name]);
			}
		}
	}

	// Prevent user input from creating / deleting saved data.
	freeze() {
		this.frozen = true;
		this.rootElem.classList.add('frozen');
	}

	private buildCreateContainer(): HTMLElement {
		let savedDataCreateFragment = document.createElement('fragment');
		savedDataCreateFragment.innerHTML = `
			<div class="saved-data-create-container">
				<label class="form-label">${this.config.label} Name</label>
				<input class="saved-data-save-input form-control" type="text" placeholder="Name">
				<button class="saved-data-save-button btn btn-primary">Save ${this.config.label}</button>
			</div>
		`;

		const saveButton = savedDataCreateFragment.querySelector('.saved-data-save-button') as HTMLButtonElement;

		saveButton.addEventListener('click', event => {
			if (this.frozen)
				return;

			const newName = this.saveInput?.value;
			if (!newName) {
				alert(`Choose a label for your saved ${this.config.label}!`);
				return;
			}

			if (newName in this.presets) {
				alert(`${this.config.label} with name ${newName} already exists.`);
				return;
			}

			this.addSavedData({
				name: newName,
				data: this.config.getData(this.modObject),
			});
			this.saveUserData();
		});

		return savedDataCreateFragment.children[0] as HTMLElement;
	}
}
