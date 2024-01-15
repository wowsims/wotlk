import { Player } from "ui/core/player";
import { IndividualSimUI } from "../../individual_sim_ui";
import {
	Potions,
	Flask,
	Food,
	Profession,
	Spec,
	Stat,
	WeaponImbue,
	Conjured
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
		this.buildFlaskPicker();
		this.buildWeaponImbuePicker();
		this.buildFoodPicker();
		this.buildPhysicalBuffPicker();
		this.buildSpellPowerBuffPicker();
		this.buildEngPicker();
		this.buildPetPicker();
	}

	private buildPotionsPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Potions</label>
        <div class="consumes-row-inputs">
          <div class="consumes-potions"></div>
          <div class="consumes-conjured"></div>
        </div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);

		const potionOptions = this.simUI.splitRelevantOptions([
			{ item: Potions.LesserManaPotion, stats: [Stat.StatIntellect] },
			{ item: Potions.ManaPotion, stats: [Stat.StatIntellect] },
		]);
		if (potionOptions.length) {
			const elem = this.rootElem.querySelector('.consumes-potions') as HTMLElement;
			new IconEnumPicker(
				elem,
				this.simUI.player,
				IconInputs.makePotionsInput(potionOptions, 'Combat Potion')
			);
		}

		const conjuredOptions = this.simUI.splitRelevantOptions([
			{ item: Conjured.ConjuredMinorRecombobulator, stats: [Stat.StatIntellect] },
			{ item: Conjured.ConjuredDemonicRune, stats: [Stat.StatIntellect] },
		]);
		if (conjuredOptions.length) {
			const elem = this.rootElem.querySelector('.consumes-conjured') as HTMLElement;
			new IconEnumPicker(elem, this.simUI.player, IconInputs.makeConjuredInput(conjuredOptions));
		}
	}

	private buildFlaskPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Elixirs</label>
        <div class="consumes-row-inputs">
          <div class="consumes-flasks"></div>
        </div>
      </div>
    `;

		const ele = this.rootElem.appendChild(fragment.children[0] as HTMLElement);

		const flaskOptions = this.simUI.splitRelevantOptions([
			{ item: Flask.FlaskOfTheTitans, stats: [Stat.StatStamina] },
			{ item: Flask.FlaskOfDistilledWisdom, stats: [Stat.StatMP5, Stat.StatSpellPower] },
			{ item: Flask.FlaskOfSupremePower, stats: [Stat.StatMP5, Stat.StatSpellPower] },
			{ item: Flask.FlaskOfChromaticResistance, stats: [Stat.StatStamina] },
		]);
		let picker: IconEnumPicker<Player<Spec>, Flask>;
		if (flaskOptions.length) {
			const elem = this.rootElem.querySelector('.consumes-flasks') as HTMLElement;
			picker = new IconEnumPicker(
				elem,
				this.simUI.player,
				IconInputs.makeFlasksInput(flaskOptions, 'Flask')
			);
		}

		// All current flasks are a level 50+ requirement
		const updateFlask = () => {
			if (this.simUI.player.getLevel() >= 50){
				picker?.restoreValue();
				ele!.classList.remove('hide');
			} else {
				picker?.storeValue();
				ele!.classList.add('hide');
			}
		};
		this.simUI.player.levelChangeEmitter.on(updateFlask);
		updateFlask();
	}

	private buildWeaponImbuePicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
    <div class="consumes-row input-root input-inline">
        <label class="form-label">Weapon Imbues</label>
        <div class="consumes-row-inputs consumes-mainhand"></div>
    </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);

		const weaponOptions = this.simUI.splitRelevantOptions([
			{ item: WeaponImbue.BrillianWizardOil, stats: [Stat.StatSpellPower] },
			{ item: WeaponImbue.BrilliantManaOil, stats: [Stat.StatHealing, Stat.StatSpellPower] },
			{ item: WeaponImbue.DenseSharpeningStone, stats: [Stat.StatAttackPower] },
			{ item: WeaponImbue.ElementalSharpeningStone, stats: [Stat.StatAttackPower] },
			{ item: WeaponImbue.BlackfathomManaOil, stats: [Stat.StatSpellPower, Stat.StatMP5] },
			{ item: WeaponImbue.BlackfathomSharpeningStone, stats: [Stat.StatMeleeHit] },
			{ item: WeaponImbue.WildStrikes, stats: [Stat.StatMeleeHit] },
		]);
		if (weaponOptions.length) {
			const elem = this.rootElem.querySelector('.consumes-mainhand') as HTMLElement;
			new IconEnumPicker(
				elem,
				this.simUI.player,	
				IconInputs.makeMainHandImbuesInput(weaponOptions, 'Weapon Imbues'),
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
			{ item: Food.FoodHotWolfRibs, stats: [Stat.StatSpirit] },
			{ item: Food.FoodSmokedSagefish, stats: [Stat.StatMP5] },
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

	private buildPhysicalBuffPicker() {
		const includeAgi = !this.simUI.individualConfig.excludeBuffDebuffInputs.includes(IconInputs.AgilityBuffInput);
		const includeStr = !this.simUI.individualConfig.excludeBuffDebuffInputs.includes(IconInputs.StrengthBuffInput);

		if (!includeAgi && !includeStr) return;

		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Physical</label>
        <div class="consumes-row-inputs consumes-physical"></div>
      </div>
    `;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const physicalConsumesElem = this.rootElem.querySelector('.consumes-physical') as HTMLElement;

		if (includeAgi){
			buildIconInput(physicalConsumesElem, this.simUI.player, IconInputs.AgilityBuffInput);
		}
		if (includeStr){
			buildIconInput(physicalConsumesElem, this.simUI.player, IconInputs.StrengthBuffInput);
		}
	}

	private buildSpellPowerBuffPicker() {
		const config = this.simUI.individualConfig;
		const includeSpellPower = config.epStats.includes(Stat.StatSpellPower) && !config.excludeBuffDebuffInputs.includes(IconInputs.SpellDamageBuff);

		if (!includeSpellPower) return;

		const includeShadowPower = !config.excludeBuffDebuffInputs.includes(IconInputs.ShadowDamageBuff);
		const includeFirePower = !config.excludeBuffDebuffInputs.includes(IconInputs.FireDamageBuff);
		const includeFrostPower = !config.excludeBuffDebuffInputs.includes(IconInputs.FrostDamageBuff);

		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Spells</label>
        <div class="consumes-row-inputs consumes-spells"></div>
      </div>
    `;

		const spellsGroup = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const spellsCnsumesElem = this.rootElem.querySelector('.consumes-spells') as HTMLElement;

		if (includeSpellPower){
			buildIconInput(spellsCnsumesElem, this.simUI.player, IconInputs.SpellDamageBuff);
		}
		if (includeFirePower){
			buildIconInput(spellsCnsumesElem, this.simUI.player, IconInputs.FireDamageBuff);
		}
		if (includeShadowPower){
			buildIconInput(spellsCnsumesElem, this.simUI.player, IconInputs.ShadowDamageBuff);
		}
		if( includeFrostPower){
			buildIconInput(spellsCnsumesElem, this.simUI.player, IconInputs.FrostDamageBuff);
		}

		const updateSpellGroup = () => {
			if (this.simUI.player.getLevel() >= 25){
				spellsGroup!.classList.remove('hide');
			} else {
				spellsGroup!.classList.add('hide');
			}
		};
		this.simUI.player.levelChangeEmitter.on(updateSpellGroup);
		updateSpellGroup();
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

		buildIconInput(tradeConsumesElem, this.simUI.player, IconInputs.Sapper);
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
