import { IndividualSimUI } from "../../individual_sim_ui";
import {
	Flask,
	Food,
	Profession,
	Spec,
	Stat
} from "../../proto/common";
import { Component } from "../component";
import { IconEnumPicker } from "../icon_enum_picker";

import * as IconInputs from '../icon_inputs.js';
import { buildIconInput } from "../icon_inputs.js";
import { SettingsTab } from "./settings_tab";

export class ConsumesPicker extends Component {
	protected settingsTab: SettingsTab;
	protected simUI: IndividualSimUI<Spec>;

	constructor(parentElem: HTMLElement, settingsTab: SettingsTab, simUI: IndividualSimUI<Spec>) {
		super(parentElem, 'consumes-picker-root');
		this.settingsTab = settingsTab;
		this.simUI = simUI;

		this.buildPotionsPicker();
		this.buildElixirsPicker();
		this.buildFoodPicker();
		this.buildEngPicker();
		this.buildPetPicker();
	}

	private buildPotionsPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Potions</label>
        <div class="consumes-row-inputs">
          <div class="consumes-prepot"></div>
          <div class="consumes-potions"></div>
          <div class="consumes-conjured"></div>
        </div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);

		// TODO: Classic - Potions aren't really a combat thing minus mana
		// const prepopPotionOptions = this.simUI.splitRelevantOptions([
		// 	// This list is smaller because some potions don't make sense to use as prepot.
		// 	// E.g. healing/mana potions.
		// 	{ item: Potions.IndestructiblePotion, stats: [Stat.StatArmor] },
		// 	{ item: Potions.InsaneStrengthPotion, stats: [Stat.StatStrength] },
		// 	{ item: Potions.HeroicPotion, stats: [Stat.StatStamina] },
		// 	{ item: Potions.PotionOfSpeed, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
		// 	{ item: Potions.PotionOfWildMagic, stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit, Stat.StatSpellPower] },
		// ]);
		// if (prepopPotionOptions.length) {
		// 	const elem = this.rootElem.querySelector('.consumes-prepot') as HTMLElement;
		// 	new IconEnumPicker(
		// 		elem,
		// 		this.simUI.player,
		// 		IconInputs.makePrepopPotionsInput(prepopPotionOptions, 'Prepop Potion (1s before combat)')
		// 	);
		// }

		// const potionOptions = this.simUI.splitRelevantOptions([
		// 	{ item: Potions.RunicHealingPotion, stats: [Stat.StatStamina] },
		// 	{ item: Potions.RunicHealingInjector, stats: [Stat.StatStamina] },
		// 	{ item: Potions.RunicManaPotion, stats: [Stat.StatIntellect] },
		// 	{ item: Potions.RunicManaInjector, stats: [Stat.StatIntellect] },
		// 	{ item: Potions.IndestructiblePotion, stats: [Stat.StatArmor] },
		// 	{ item: Potions.InsaneStrengthPotion, stats: [Stat.StatStrength] },
		// 	{ item: Potions.HeroicPotion, stats: [Stat.StatStamina] },
		// 	{ item: Potions.PotionOfSpeed, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
		// 	{ item: Potions.PotionOfWildMagic, stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit, Stat.StatSpellPower] },
		// ]);
		// if (potionOptions.length) {
		// 	const elem = this.rootElem.querySelector('.consumes-potions') as HTMLElement;
		// 	new IconEnumPicker(
		// 		elem,
		// 		this.simUI.player,
		// 		IconInputs.makePotionsInput(potionOptions, 'Combat Potion')
		// 	);
		// }

		// TODO: Classic Use APL?
		// const conjuredOptions = this.simUI.splitRelevantOptions([
		// 	this.simUI.player.getClass() == Class.ClassRogue ? { item: Conjured.ConjuredRogueThistleTea, stats: [] } : null,
		// 	{ item: Conjured.ConjuredHealthstone, stats: [Stat.StatStamina] },
		// 	{ item: Conjured.ConjuredDarkRune, stats: [Stat.StatIntellect] },
		// ]);
		// if (conjuredOptions.length) {
		// 	const elem = this.rootElem.querySelector('.consumes-conjured') as HTMLElement;
		// 	new IconEnumPicker(elem, this.simUI.player, IconInputs.makeConjuredInput(conjuredOptions));
		// }
	}

	private buildElixirsPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Elixirs</label>
        <div class="consumes-row-inputs">
          <div class="consumes-flasks"></div>
        </div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);

		const flaskOptions = this.simUI.splitRelevantOptions([
			{ item: Flask.FlaskOfTheTitans, stats: [Stat.StatStamina] },
			{ item: Flask.FlaskOfDistilledWisdom, stats: [Stat.StatMP5, Stat.StatSpellPower] },
			{ item: Flask.FlaskOfSupremePower, stats: [Stat.StatMP5, Stat.StatSpellPower] },
			{ item: Flask.FlaskOfChromaticResistance, stats: [Stat.StatStamina] },
		]);
		if (flaskOptions.length) {
			const elem = this.rootElem.querySelector('.consumes-flasks') as HTMLElement;
			new IconEnumPicker(
				elem,
				this.simUI.player,
				IconInputs.makeFlasksInput(flaskOptions, 'Flask')
			);
		}
	}

	private buildFoodPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Food</label>
        <div class="consumes-row-inputs">
          <div class="consumes-food"></div>
        </div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);

		const foodOptions = this.simUI.splitRelevantOptions([
			{ item: Food.FoodNightfinSoup, stats: [Stat.StatMP5, Stat.StatSpellPower] },
			{ item: Food.FoodGrilledSquid, stats: [Stat.StatAgility] },
			{ item: Food.FoodSmokedDesertDumpling, stats: [Stat.StatStrength] },
			{ item: Food.FoodRunnTumTuberSurprise, stats: [Stat.StatIntellect] },
			{ item: Food.FoodDirgesKickChimaerokChops, stats: [Stat.StatStamina] },
			{ item: Food.FoodBlessSunfruit, stats: [Stat.StatStrength] },
			{ item: Food.FoodBlessedSunfruitJuice, stats: [Stat.StatSpirit] },
		]);
		if (foodOptions.length) {
			const elem = this.rootElem.querySelector('.consumes-food') as HTMLElement;
			new IconEnumPicker(elem, this.simUI.player, IconInputs.makeFoodInput(foodOptions));
		}
	}

	private buildEngPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Engineering</label>
        <div class="consumes-row-inputs consumes-trade"></div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);

		const tradeConsumesElem = this.rootElem.querySelector('.consumes-trade') as HTMLElement;

		// TODO Classic
		// buildIconInput(tradeConsumesElem, this.simUI.player, IconInputs.Sapper);
		buildIconInput(tradeConsumesElem, this.simUI.player, IconInputs.FillerExplosiveInput);

		const updateProfession = () => {
			if (this.simUI.player.hasProfession(Profession.Engineering))
				tradeConsumesElem.parentElement!.classList.remove('hide');
			else
				tradeConsumesElem.parentElement!.classList.add('hide');
		};
		this.simUI.player.professionChangeEmitter.on(updateProfession);
		updateProfession();
	}

	private buildPetPicker() {
		if (this.simUI.individualConfig.petConsumeInputs?.length) {
			let fragment = document.createElement('fragment');
			fragment.innerHTML = `
        <div class="consumes-row input-root input-inline">
          <label class="form-label">Pet</label>
          <div class="consumes-row-inputs consumes-pet"></div>
        </div>
      `;

			this.rootElem.appendChild(fragment.children[0] as HTMLElement);

			const petConsumesElem = this.rootElem.querySelector('.consumes-pet') as HTMLElement;
			this.simUI.individualConfig.petConsumeInputs.map(iconInput => buildIconInput(petConsumesElem, this.simUI.player, iconInput));
		}
	}
}
