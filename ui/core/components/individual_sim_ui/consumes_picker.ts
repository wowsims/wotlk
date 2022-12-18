import { IndividualSimUI } from "../../individual_sim_ui";
import {
  BattleElixir,
  Class,
  Conjured,
  Flask,
  Food,
  GuardianElixir,
  Potions,
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

    const prepopPotionOptions = this.simUI.splitRelevantOptions([
			// This list is smaller because some potions don't make sense to use as prepot.
			// E.g. healing/mana potions.
			{ item: Potions.IndestructiblePotion, stats: [Stat.StatArmor] },
			{ item: Potions.InsaneStrengthPotion, stats: [Stat.StatStrength] },
			{ item: Potions.HeroicPotion, stats: [Stat.StatStamina] },
			{ item: Potions.PotionOfSpeed, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
			{ item: Potions.PotionOfWildMagic, stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit, Stat.StatSpellPower] },
		]);
		if (prepopPotionOptions.length) {
			const elem = this.rootElem.querySelector('.consumes-prepot') as HTMLElement;
			new IconEnumPicker(
        elem,
        this.simUI.player,
        IconInputs.makePrepopPotionsInput(prepopPotionOptions, 'Prepop Potion (1s before combat)')
      );
    }

		const potionOptions = this.simUI.splitRelevantOptions([
			{ item: Potions.RunicHealingPotion, stats: [Stat.StatStamina] },
			{ item: Potions.RunicHealingInjector, stats: [Stat.StatStamina] },
			{ item: Potions.RunicManaPotion, stats: [Stat.StatIntellect] },
			{ item: Potions.RunicManaInjector, stats: [Stat.StatIntellect] },
			{ item: Potions.IndestructiblePotion, stats: [Stat.StatArmor] },
			{ item: Potions.InsaneStrengthPotion, stats: [Stat.StatStrength] },
			{ item: Potions.HeroicPotion, stats: [Stat.StatStamina] },
			{ item: Potions.PotionOfSpeed, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
			{ item: Potions.PotionOfWildMagic, stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit, Stat.StatSpellPower] },
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
			this.simUI.player.getClass() == Class.ClassRogue ? { item: Conjured.ConjuredRogueThistleTea, stats: [] } : null,
			{ item: Conjured.ConjuredHealthstone, stats: [Stat.StatStamina] },
			{ item: Conjured.ConjuredDarkRune, stats: [Stat.StatIntellect] },
			{ item: Conjured.ConjuredFlameCap, stats: [] },
		]);
		if (conjuredOptions.length) {
			const elem = this.rootElem.querySelector('.consumes-conjured') as HTMLElement;
			new IconEnumPicker(elem, this.simUI.player, IconInputs.makeConjuredInput(conjuredOptions));
		}
  }

  private buildElixirsPicker() {
    let fragment = document.createElement('fragment');
    fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Elixirs</label>
        <div class="consumes-row-inputs">
          <div class="consumes-flasks"></div>
          <span class="elixir-space">or</span>
          <div class="consumes-battle-elixirs"></div>
          <div class="consumes-guardian-elixirs"></div>
        </div>
      </div>
    `;

    this.rootElem.appendChild(fragment.children[0] as HTMLElement);

    const flaskOptions = this.simUI.splitRelevantOptions([
			{ item: Flask.FlaskOfTheFrostWyrm, stats: [Stat.StatSpellPower] },
			{ item: Flask.FlaskOfEndlessRage, stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower] },
			{ item: Flask.FlaskOfPureMojo, stats: [Stat.StatMP5] },
			{ item: Flask.FlaskOfStoneblood, stats: [Stat.StatStamina] },
			{ item: Flask.LesserFlaskOfToughness, stats: [Stat.StatResilience] },
			{ item: Flask.LesserFlaskOfResistance, stats: [Stat.StatArcaneResistance, Stat.StatFireResistance, Stat.StatFrostResistance, Stat.StatNatureResistance, Stat.StatShadowResistance] },
		]);
		if (flaskOptions.length) {
			const elem = this.rootElem.querySelector('.consumes-flasks') as HTMLElement;
			new IconEnumPicker(
        elem,
        this.simUI.player,
        IconInputs.makeFlasksInput(flaskOptions, 'Flask')
      );
		}

		const battleElixirOptions = this.simUI.splitRelevantOptions([
			{ item: BattleElixir.ElixirOfAccuracy, stats: [Stat.StatMeleeHit, Stat.StatSpellHit] },
			{ item: BattleElixir.ElixirOfArmorPiercing, stats: [Stat.StatArmorPenetration] },
			{ item: BattleElixir.ElixirOfDeadlyStrikes, stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit] },
			{ item: BattleElixir.ElixirOfExpertise, stats: [Stat.StatExpertise] },
			{ item: BattleElixir.ElixirOfLightningSpeed, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
			{ item: BattleElixir.ElixirOfMightyAgility, stats: [Stat.StatAgility] },
			{ item: BattleElixir.ElixirOfMightyStrength, stats: [Stat.StatStrength] },
			{ item: BattleElixir.GurusElixir, stats: [Stat.StatStamina, Stat.StatAgility, Stat.StatStrength, Stat.StatSpirit, Stat.StatIntellect] },
			{ item: BattleElixir.SpellpowerElixir, stats: [Stat.StatSpellPower] },
			{ item: BattleElixir.WrathElixir, stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower] },
		]);

    const battleElixirsContainer = this.rootElem.querySelector('.consumes-battle-elixirs') as HTMLElement;
		if (battleElixirOptions.length) {
			new IconEnumPicker(
        battleElixirsContainer,
        this.simUI.player,
        IconInputs.makeBattleElixirsInput(battleElixirOptions, 'Battle Elixir')
      );
		} else {
      battleElixirsContainer.remove();
    }

		const guardianElixirOptions = this.simUI.splitRelevantOptions([
			{ item: GuardianElixir.ElixirOfMightyDefense, stats: [Stat.StatDefense] },
			{ item: GuardianElixir.ElixirOfMightyFortitude, stats: [Stat.StatStamina] },
			{ item: GuardianElixir.ElixirOfMightyMageblood, stats: [Stat.StatMP5] },
			{ item: GuardianElixir.ElixirOfMightyThoughts, stats: [Stat.StatIntellect] },
			{ item: GuardianElixir.ElixirOfProtection, stats: [Stat.StatArmor] },
			{ item: GuardianElixir.ElixirOfSpirit, stats: [Stat.StatSpirit] },
			{ item: GuardianElixir.GiftOfArthas, stats: [Stat.StatStamina] },
		]);

    const guardianElixirsContainer = this.rootElem.querySelector('.consumes-guardian-elixirs') as HTMLElement;
		if (guardianElixirOptions.length) {
			const guardianElixirsContainer = this.rootElem.querySelector('.consumes-guardian-elixirs') as HTMLElement;
			new IconEnumPicker(
        guardianElixirsContainer,
        this.simUI.player,
        IconInputs.makeGuardianElixirsInput(guardianElixirOptions, 'Guardian Elixir')
      );
		} else {
      guardianElixirsContainer.remove();
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
			{ item: Food.FoodFishFeast, stats: [Stat.StatStamina, Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower] },
			{ item: Food.FoodGreatFeast, stats: [Stat.StatStamina, Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower] },
			{ item: Food.FoodBlackenedDragonfin, stats: [Stat.StatAgility] },
			{ item: Food.FoodDragonfinFilet, stats: [Stat.StatStrength] },
			{ item: Food.FoodCuttlesteak, stats: [Stat.StatSpirit] },
			{ item: Food.FoodMegaMammothMeal, stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower] },
			{ item: Food.FoodHeartyRhino, stats: [Stat.StatArmorPenetration] },
			{ item: Food.FoodRhinoliciousWormsteak, stats: [Stat.StatExpertise] },
			{ item: Food.FoodFirecrackerSalmon, stats: [Stat.StatSpellPower] },
			{ item: Food.FoodSnapperExtreme, stats: [Stat.StatMeleeHit, Stat.StatSpellHit] },
			{ item: Food.FoodSpicedWormBurger, stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit] },
			{ item: Food.FoodImperialMantaSteak, stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
			{ item: Food.FoodMightyRhinoDogs, stats: [Stat.StatMP5] },
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

		buildIconInput(tradeConsumesElem, this.simUI.player, IconInputs.ThermalSapper);
		buildIconInput(tradeConsumesElem, this.simUI.player, IconInputs.ExplosiveDecoy);
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
