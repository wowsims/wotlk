import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { Debuffs, IndividualBuffs, PartyBuffs, PseudoStat, RaidBuffs, Spec, Stat, TristateEffect } from '../core/proto/common.js';
import { Gear } from '../core/proto_utils/gear.js';
import { Stats } from '../core/proto_utils/stats.js';
import { TypedEvent } from '../core/typed_event.js';

import * as OtherInputs from '../core/components/other_inputs.js';

import * as WarriorInputs from './inputs.js';
import * as Presets from './presets.js';

export class WarriorSimUI extends IndividualSimUI<Spec.SpecWarrior> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecWarrior>) {
		super(parentElem, player, {
			cssClass: 'warrior-sim-ui',
			cssScheme: 'warrior',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatExpertise,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatArmor,
			],
			epPseudoStats: [
				PseudoStat.PseudoStatMainHandDps,
				PseudoStat.PseudoStatOffHandDps,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatStamina,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatExpertise,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatArmor,
			],
			modifyDisplayStats: (_: Player<Spec.SpecWarrior>) => {
				let stats = new Stats();

				return {
					talents: stats,
				};
			},

			defaults: {
				// Default equipped gear.
				gear: Presets.EMPTY_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatStrength]: 2.72,
					[Stat.StatAgility]: 1.82,
					[Stat.StatAttackPower]: 1,
					[Stat.StatExpertise]: 2.55,
					[Stat.StatMeleeHit]: 0.79,
					[Stat.StatMeleeCrit]: 2.12,
					[Stat.StatMeleeHaste]: 1.72,
					[Stat.StatArmorPenetration]: 2.17,
					[Stat.StatArmor]: 0.03,
				}, {
					[PseudoStat.PseudoStatMainHandDps]: 6.29,
					[PseudoStat.PseudoStatOffHandDps]: 3.58,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default talents.
				talents: Presets.Talent25.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					leaderOfThePack: true,
					devotionAura: TristateEffect.TristateEffectImproved,
					stoneskinTotem: TristateEffect.TristateEffectImproved,
				}),
				partyBuffs: PartyBuffs.create({
					heroicPresence: false,
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
				}),
				debuffs: Debuffs.create({
					sunderArmor: true,
					curseOfWeakness: TristateEffect.TristateEffectRegular,
					faerieFire: TristateEffect.TristateEffectImproved,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				WarriorInputs.ShoutPicker,
				WarriorInputs.Recklessness,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: WarriorInputs.WarriorRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					WarriorInputs.StartingRage,
					WarriorInputs.StanceSnapshot,
					OtherInputs.TankAssignment,
					OtherInputs.InFrontOfTarget,
				],
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: true,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.Talent25
				],
				// Preset rotations that the user can quickly select.
				rotations: [
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
				],
			},
		});
	}

	calcExpCap(): Stats {
		let expCap = 6.5 * 32.79 + 4;

		return new Stats().withStat(Stat.StatExpertise, expCap);
	}

	calcArpCap(gear: Gear): Stats {
		let arpCap = 1404;

		if (gear.hasTrinket(45931)) {
			arpCap = 659;
		} else if (gear.hasTrinket(40256)) {
			arpCap = 798;
		}

		return new Stats().withStat(Stat.StatArmorPenetration, arpCap);
	}

	calcArpTarget(gear: Gear): number {
		if (gear.hasTrinket(45931)) {
			return 648;
		}

		if (gear.hasTrinket(40256)) {
			return 787;
		}

		return 1399;
	}

	calcCritCap(gear: Gear): Stats {
		const baseCritCapPercentage = 77.8; // includes 3% Crit debuff
		let agiProcs = 0;

		if (gear.hasRelic(47668)) {
			agiProcs += 200;
		}

		if (gear.hasRelic(50456)) {
			agiProcs += 44*5;
		}

		if (gear.hasTrinket(47131) || gear.hasTrinket(47464)) {
			agiProcs += 510;
		}

		if (gear.hasTrinket(47115) || gear.hasTrinket(47303)) {
			agiProcs += 450;
		}

		if (gear.hasTrinket(44253) || gear.hasTrinket(42987)) {
			agiProcs += 300;
		}

		return new Stats().withStat(Stat.StatMeleeCrit, (baseCritCapPercentage - agiProcs*1.1*1.06*1.02/83.33) * 45.91);
	}

	async updateGear(gear: Gear): Promise<Stats> {
		this.player.setGear(TypedEvent.nextEventID(), gear);
		await this.sim.updateCharacterStats(TypedEvent.nextEventID());
		return Stats.fromProto(this.player.getCurrentStats().finalStats);
	}
}
