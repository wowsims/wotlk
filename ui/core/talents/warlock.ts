import { WarlockTalents, WarlockMajorGlyph, WarlockMinorGlyph } from '../proto/warlock.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import WarlockTalentJson from './trees/warlock.json';

export const warlockTalentsConfig: TalentsConfig<WarlockTalents> = newTalentsConfig(WarlockTalentJson);

export const warlockGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[WarlockMajorGlyph.GlyphOfChaosBolt]: {
			name: 'Glyph of Chaos Bolt',
			description: 'Reduces the cooldown on Chaos Bolt by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warlock_chaosbolt.jpg',
		},
		[WarlockMajorGlyph.GlyphOfConflagrate]: {
			name: 'Glyph of Conflagrate',
			description: 'Your Conflagrate spell no longer consumes your Immolate or Shadowflame spell from the target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_fireball.jpg',
		},
		[WarlockMajorGlyph.GlyphOfCorruption]: {
			name: 'Glyph of Corruption',
			description: 'Your Corruption spell has a 4% chance to cause you to enter a Shadow Trance state after damaging the opponent. The Shadow Trance state reduces the casting time of your next Shadow Bolt spell by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_abominationexplosion.jpg',
		},
		[WarlockMajorGlyph.GlyphOfCurseOfAgony]: {
			name: 'Glyph of Curse of Agony',
			description: 'Increases the duration of your Curse of Agony by 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_curseofsargeras.jpg',
		},
		[WarlockMajorGlyph.GlyphOfDeathCoil]: {
			name: 'Glyph of Death Coil',
			description: 'Increases the duration of your Death Coil by 0.5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathcoil.jpg',
		},
		[WarlockMajorGlyph.GlyphOfDemonicCircle]: {
			name: 'Glyph of Demonic Circle',
			description: 'Reduces the cooldown on Demonic Circle by 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demoniccirclesummon.jpg',
		},
		[WarlockMajorGlyph.GlyphOfFear]: {
			name: 'Glyph of Fear',
			description: 'Increases the damage your Fear target can take before the Fear effect is removed by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_possession.jpg',
		},
		[WarlockMajorGlyph.GlyphOfFelguard]: {
			name: 'Glyph of Felguard',
			description: 'Increases the Felguard\'s total attack power by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonfelguard.jpg',
		},
		[WarlockMajorGlyph.GlyphOfFelhunter]: {
			name: 'Glyph of Felhunter',
			description: 'When your Felhunter uses Devour Magic, you will also be healed for that amount.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonfelhunter.jpg',
		},
		[WarlockMajorGlyph.GlyphOfHaunt]: {
			name: 'Glyph of Haunt',
			description: 'The bonus damage granted by your Haunt spell is increased by an additional 3%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warlock_haunt.jpg',
		},
		[WarlockMajorGlyph.GlyphOfHealthFunnel]: {
			name: 'Glyph of Health Funnel',
			description: 'Reduces the pushback suffered from damaging attacks while channeling your Health Funnel spell by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_lifedrain.jpg',
		},
		[WarlockMajorGlyph.GlyphOfHealthstone]: {
			name: 'Glyph of Healthstone',
			description: 'You receive 30% more healing from using a healthstone.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_stone_04.jpg',
		},
		[WarlockMajorGlyph.GlyphOfHowlOfTerror]: {
			name: 'Glyph of Howl of Terror',
			description: 'Reduces the cooldown on your Howl of Terror spell by 8 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathscream.jpg',
		},
		[WarlockMajorGlyph.GlyphOfImmolate]: {
			name: 'Glyph of Immolate',
			description: 'Increases the periodic damage of your Immolate by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_immolation.jpg',
		},
		[WarlockMajorGlyph.GlyphOfImp]: {
			name: 'Glyph of Imp',
			description: 'Increases the damage done by your Imp\'s Firebolt spell by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonimp.jpg',
		},
		[WarlockMajorGlyph.GlyphOfIncinerate]: {
			name: 'Glyph of Incinerate',
			description: 'Increases the damage done by Incinerate by 5%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_burnout.jpg',
		},
		[WarlockMajorGlyph.GlyphOfLifeTap]: {
			name: 'Glyph of Life Tap',
			description: 'When you use Life Tap or Dark Pact, you gain 20% of your Spirit as spell power for 40 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_burningspirit.jpg',
		},
		[WarlockMajorGlyph.GlyphOfMetamorphosis]: {
			name: 'Glyph of Metamorphosis',
			description: 'Increases the duration of your Metamorphosis by 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonform.jpg',
		},
		[WarlockMajorGlyph.GlyphOfQuickDecay]: {
			name: 'Glyph of Quick Decay',
			description: 'Your haste now reduces the time between periodic damage ticks of your Corruption spell.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_abominationexplosion.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSearingPain]: {
			name: 'Glyph of Searing Pain',
			description: 'Increases the critical strike bonus of your Searing Pain by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_soulburn.jpg',
		},
		[WarlockMajorGlyph.GlyphOfShadowBolt]: {
			name: 'Glyph of Shadow Bolt',
			description: 'Reduces the mana cost of your Shadow Bolt by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowbolt.jpg',
		},
		[WarlockMajorGlyph.GlyphOfShadowburn]: {
			name: 'Glyph of Shadowburn',
			description: 'Increases the critical strike chance of Shadowburn by 20% when the target is below 35% health.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_scourgebuild.jpg',
		},
		[WarlockMajorGlyph.GlyphOfShadowflame]: {
			name: 'Glyph of Shadowflame',
			description: 'Your Shadowflame also applies a 70% movement speed slow on its victims.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warlock_shadowflame.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSiphonLife]: {
			name: 'Glyph of Siphon Life',
			description: 'Increases the healing you receive from your Siphon Life talent by 25%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_requiem.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSoulLink]: {
			name: 'Glyph of Soul Link',
			description: 'Increases the percentage of damage shared via your Soul Link by an additional 5%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_gathershadows.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSoulstone]: {
			name: 'Glyph of Soulstone',
			description: 'Increases the amount of health you gain from resurrecting via a Soulstone by 300%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_soulgem.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSuccubus]: {
			name: 'Glyph of Succubus',
			description: 'Your Succubus\'s Seduction ability also removes all damage over time effects from the target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonsuccubus.jpg',
		},
		[WarlockMajorGlyph.GlyphOfUnstableAffliction]: {
			name: 'Glyph of Unstable Affliction',
			description: 'Decreases the casting time of your Unstable Affliction by 0.2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_unstableaffliction_3.jpg',
		},
		[WarlockMajorGlyph.GlyphOfVoidwalker]: {
			name: 'Glyph of Voidwalker',
			description: 'Increases your Voidwalker\'s total Stamina by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonvoidwalker.jpg',
		},
	},
	minorGlyphs: {
		[WarlockMinorGlyph.GlyphOfCurseOfExhausion]: {
			name: 'Glyph of Curse of Exhausion',
			description: 'Increases the range of your Curse of Exhaustion spell by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_grimward.jpg',
		},
		[WarlockMinorGlyph.GlyphOfDrainSoul]: {
			name: 'Glyph of Drain Soul',
			description: 'Your Drain Soul ability occasionally creates an additional soul shard.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_haunting.jpg',
		},
		[WarlockMinorGlyph.GlyphOfSubjugateDemon]: {
			name: 'Glyph of Subjugate Demon',
			description: 'Reduces the cast time of your Subjugate Demon spell by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_enslavedemon.jpg',
		},
		[WarlockMinorGlyph.GlyphOfKilrogg]: {
			name: 'Glyph of Kilrogg',
			description: 'Increases the movement speed of your Eye of Kilrogg by 50% and allows it to fly in areas where flying mounts are enabled.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_evileye.jpg',
		},
		[WarlockMinorGlyph.GlyphOfSouls]: {
			name: 'Glyph of Souls',
			description: 'Reduces the mana cost of your Ritual of Souls spell by 70%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadesofdarkness.jpg',
		},
		[WarlockMinorGlyph.GlyphOfUnendingBreath]: {
			name: 'Glyph of Unending Breath',
			description: 'Increases the swim speed of targets affected by your Unending Breath spell by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonbreath.jpg',
		},
	},
};
