import { Exporter } from '../core/components/exporters';
import { Importer } from '../core/components/importers';
import { RaidSimSettings } from '../core/proto/ui';
import { EventID, TypedEvent } from '../core/typed_event';
import { Party as PartyProto, Player as PlayerProto, Raid as RaidProto } from '../core/proto/api';
import {
	Class,
	Encounter as EncounterProto,
	EquipmentSpec,
	Faction,
	ItemSpec,
	Profession,
	Race,
	UnitReference,
	Spec,
	Target as TargetProto,
	UnitReference_Type,
} from '../core/proto/common';
import { professionNames, raceNames } from '../core/proto_utils/names';
import {
	DruidSpecs,
	DeathknightSpecs,
	PriestSpecs,
	RogueSpecs,
	SpecOptions,
	getTalentTreePoints,
	makeDefaultBlessings,
	raceToFaction,
	isTankSpec,
	playerToSpec,
} from '../core/proto_utils/utils';
import { MAX_NUM_PARTIES } from '../core/raid';
import { RaidSimPreset } from '../core/individual_sim_ui';
import { Player } from '../core/player';
import { Encounter } from '../core/encounter';
import { bucket, distinct } from '../core/utils';

import { playerPresets } from './presets';
import { RaidSimUI } from './raid_sim_ui';

export class RaidJsonImporter extends Importer {
	private readonly simUI: RaidSimUI;
	constructor(parent: HTMLElement, simUI: RaidSimUI) {
		super(parent, simUI, 'JSON Import', true);
		this.simUI = simUI;

		this.descriptionElem.innerHTML = `
			<p>
				Import settings from a JSON text file, which can be created using the JSON Export feature of this site.
			</p>
			<p>
				To import, paste the JSON text below and click, 'Import'.
			</p>
		`;
	}

	onImport(data: string) {
		const settings = RaidSimSettings.fromJsonString(data, { ignoreUnknownFields: true });
		this.simUI.fromProto(TypedEvent.nextEventID(), settings);
		this.close();
	}
}

export class RaidJsonExporter extends Exporter {
	private readonly simUI: RaidSimUI;

	constructor(parent: HTMLElement, simUI: RaidSimUI) {
		super(parent, simUI, {title: 'JSON Export', allowDownload: true});
		this.simUI = simUI;
		this.init();
	}

	getData(): string {
		return JSON.stringify(RaidSimSettings.toJson(this.simUI.toProto()), null, 2);
	}
}

export class RaidWCLImporter extends Importer {

	private queryCounter: number = 0;

	private readonly simUI: RaidSimUI;
	constructor(parent: HTMLElement, simUI: RaidSimUI) {
		super(parent, simUI, 'WCL Import', false);
		this.simUI = simUI;
		this.textElem.classList.add('small-textarea');
		this.descriptionElem.innerHTML = `
			<p>
				Imports the entire raid from a WCL report.<br>
			</p>
			<p>
				To import, paste the WCL report and fight link (https://classic.warcraftlogs.com/reports/REPORTID#fight=FIGHTID).<br>
				Include the fight ID or else the first fight in the report will be used.<br>
			</p>
			<p>
				The following are imported directly from the report:
				<ul>
					<li>Player Name</li>
					<li>Equipment (items, enchants, and gems)</li>
					<li>Faction (Alliance / Horde)</li>
					<li>Encounter: If the import link has a fight ID we try to match with a preset Encounter. Note that many Encounters are still unimplemented.</li>
				</ul>

				The following are not available directly from the report data, but we try to infer them:
				<ul>
					<li>Talents: Log data only gives us the tree summary (e.g. '51/20/0') so we match this with the closest preset talent build.</li>
					<li>Glyphs: Glyphs are absent from log data, but we pair them with the inferred Talents.</li>
					<li>Race: Inferred from Race-specific abilties used in any fight if possible, or defaults to Spec-specific Race.</li>
					<li>Professions: Inferred from profession-locked items/enchants/gems.</li>
					<li>Buff assignments (Innervate, Unholy Frenzy, etc): Inferred from casts.</li>
					<li>Party Composition: Inferred from party-only effects, such as Heroic Presence, Prayer of Healing, or Vampiric Touch.</li>
				</ul>

				The following are not imported, and instead use spec-specific defaults:
				<ul>
					<li>Rotation / Spec-specific options</li>
					<li>Consumes</li>
					<li>Paladin Blessings</li>
				</ul>
			</p>
		`;
	}

	private token: string = '';
	private async getWCLBearerToken(): Promise<string> {
		if (this.token == '') {
			const response = await fetch('https://classic.warcraftlogs.com/oauth/token', {
				'method': 'POST',
				'headers': {
					'Authorization': 'Basic ' + btoa('963d31c8-7efa-4dde-87cf-1b254a8a2f8c:lRJVhujEEnF96xfUoxVHSpnqKN9v8bTqGEjutsO3'),
				},
				body: new URLSearchParams({
					'grant_type': 'client_credentials',
				}),
			})
			const json = await response.json();
			this.token = json.access_token;
		}
		return this.token;
	}

	private async queryWCL(query: string): Promise<any> {
		const token = await this.getWCLBearerToken();
		const headers = {
			'Content-Type': 'application/json',
			'Authorization': `Bearer ${token}`,
			'Accept': 'application/json',
		};

		const queryURL = `https://classic.warcraftlogs.com/api/v2/client?query=${query}`;
		this.queryCounter++;

		// Query WCL
		const res = await fetch(encodeURI(queryURL), {
			'method': 'GET',
			'headers': headers,
		});

		const result = await res.json();
		if (result?.errors?.length) {
			const errorStr = result.errors.map((e: any) => e.message).join('\n');
			throw new Error(`GraphQL error: ${errorStr}\n\nQuery: ${query}`);
		} else {
			console.debug(`WCL query: ${query}\n\nResult: ${JSON.stringify(result)}`);
		}
		return result;
	}

	private async parseURL(url: string): Promise<wclUrlData> {
		const match = url.match(/classic\.warcraftlogs\.com\/reports\/([a-zA-Z0-9:]+)\/?(#.*fight=((\d+)|(last)))?/);
		if (!match) {
			throw new Error(`Invalid WCL URL ${url}, must look like "classic.warcraftlogs.com/reports/XXXX"`);
		}

		const urlData = {
			reportID: match[1],
			fightID: '',
		}

		// If the URL has a Fight ID in it, use it
		if (match[2] && match[3] && match[3] != 'last') {
			urlData.fightID = match[3];
		} else {
			// Make a separate query to get the corresponding ReportFights
			const fightDataQuery = `{
				reportData {
					report(code: "${urlData.reportID}") {
						fights(killType: Kills, translate: true) {
							id, name
						}
					}
				}
			}`;

			const fightData = await this.queryWCL(fightDataQuery);
			const fights = fightData.data.reportData.report.fights;

			if (match[3] == 'last') {
				urlData.fightID = String(fights[fights.length - 1].id)
			} else {
				// Default to using the first Fight
				urlData.fightID = String(fights[0].id);
			}
		}

		console.debug(`Importing WCL report: ${JSON.stringify(urlData)}`);
		return urlData;
	}

	private async getRateLimit(): Promise<wclRateLimitData> {
		const query = `{
	    rateLimitData {
	      limitPerHour, pointsSpentThisHour, pointsResetIn
	    }
	  }`;
		const result = await this.queryWCL(query);
		const data = result['data']['rateLimitData'] as wclRateLimitData;
		return data;
	}

	async onImport(importLink: string) {
		this.importButton.disabled = true;
		this.rootElem.style.cursor = 'wait';
		try {
			await this.doImport(importLink);
		} catch (error) {
			console.error(error);
			alert('Failed import from WCL: ' + error);
		}
		this.importButton.disabled = false
		this.rootElem.style.removeProperty('cursor');
	}

	async doImport(importLink: string) {
		if (!importLink.length) {
			throw new Error('No import link provided!');
		}

		const urlData = await this.parseURL(importLink);
		const _rateLimit = await this.getRateLimit();

		// Schema for WCL API here: https://www.warcraftlogs.com/v2-api-docs/warcraft/
		// WCL charges us 1 'point' for each subquery we issue within the request. So
		// by using filter expressions we can combine our queries together, to spend
		// fewer points.
		const reportDataQuery = `{
			reportData {
				report(code: "${urlData.reportID}") {
					guild {
						name faction {id}
					}
					playerDetails: table(fightIDs: [${urlData.fightID}], dataType: Casts, killType: All, viewBy: Default)
					combatantInfoEvents: events(fightIDs: [${urlData.fightID}], dataType:CombatantInfo, limit: 50) { data }
					fights(fightIDs: [${urlData.fightID}]) {
						startTime, endTime, id, name
					}

					reportCastEvents: events(dataType:Casts, endTime: 99999999, filterExpression: "${[racialSpells, professionSpells].flat().map(spell => spell.id).map(id => `ability.id = ${id}`).join(' OR ')
			}", limit: 10000) { data }

					fightCastEvents: events(fightIDs: [${urlData.fightID}], dataType:Casts, filterExpression: "${[externalCDSpells].flat().map(spell => spell.id).map(id => `ability.id = ${id}`).join(' OR ')
			}", limit: 10000) { data }

					fightHealEvents: events(fightIDs: [${urlData.fightID}], dataType:Healing, filterExpression: "${[samePartyHealingSpells, otherPartyHealingSpells].flat().map(spell => spell.id).map(id => `ability.id = ${id}`).join(' OR ')
			}", limit: 10000) { data }

					manaTideTotem: events(fightIDs: [${urlData.fightID}], dataType:Resources, filterExpression: "ability.id = 39609", limit: 100) { data }
				}
			}
		}`;
		const reportData = await this.queryWCL(reportDataQuery);

		// Process the report data.
		const wclData = reportData.data.reportData.report; // TODO: Typings?
		const playerData: wclPlayer[] = wclData.playerDetails.data.entries;

		TypedEvent.freezeAllAndDo(() => {
			const eventID = TypedEvent.nextEventID();
			const wclPlayers = playerData.map(wclPlayer => new WCLSimPlayer(wclPlayer, this.simUI, eventID));
			this.inferRace(eventID, wclData, wclPlayers);
			this.inferProfessions(eventID, wclData, wclPlayers);
			this.inferAssignments(eventID, wclData, wclPlayers);
			this.inferPartyComposition(eventID, wclData, wclPlayers);
			const numPaladins = wclPlayers.filter(player => player.player.getClass() == Class.ClassPaladin).length;
			const settings = RaidSimSettings.create({
				encounter: this.getEncounterProto(wclData),
				raid: this.getRaidProto(wclPlayers),
				blessings: makeDefaultBlessings(numPaladins),
			});

			// Clear the raid out to avoid any taint issues.
			this.simUI.clearRaid(eventID);
			this.simUI.fromProto(eventID, settings);
		});

		this.close();
	}

	private inferRace(eventID: EventID, wclData: any, wclPlayers: WCLSimPlayer[]) {
		wclPlayers.forEach(p => p.player.setRace(eventID, Race.RaceUnknown));

		// If defined in log, use that faction. Otherwise default to UI setting.
		let faction = (wclData.guild?.faction?.id || this.simUI.raidPicker?.getCurrentFaction() || Faction.Horde) as Faction;

		wclData.combatantInfoEvents.data.forEach((combatantInfo: wclCombatantInfoEvent) => {
			combatantInfo.auras
				.filter(aura => aura.ability == 28878)
				.forEach(aura => {
					const sourcePlayer = wclPlayers.find(player => player.id == aura.source);
					if (sourcePlayer && sourcePlayer.player.getRace() != Race.RaceDraenei) {
						console.log(`Inferring player ${sourcePlayer.name} has race ${raceNames.get(Race.RaceDraenei)} from Heroic Presence aura event`);
						sourcePlayer.player.setRace(eventID, Race.RaceDraenei);
						faction = Faction.Alliance;
					}
				});
		});

		const castEventsBySpellId = bucket(wclData.reportCastEvents.data as Array<wclCastEvent>, event => String(event.abilityGameID));
		racialSpells.forEach(spell => {
			const spellEvents: Array<wclCastEvent> = castEventsBySpellId[spell.id] || [];
			spellEvents.forEach(event => {
				const sourcePlayer = wclPlayers.find(player => player.id == event.sourceID);
				if (sourcePlayer) {
					console.log(`Inferring player ${sourcePlayer.name} has race ${raceNames.get(spell.race)} from ${spell.name} event`);
					sourcePlayer.player.setRace(eventID, spell.race);
					faction = raceToFaction[spell.race];
				}
			});
		});

		wclPlayers.forEach(p => {
			if (p.player.getRace() == Race.RaceUnknown) {
				p.player.setRace(eventID, p.preset.defaultFactionRaces[faction]);
			}
		});
	}

	private inferProfessions(eventID: EventID, wclData: any, wclPlayers: WCLSimPlayer[]) {
		const castEventsBySpellId = bucket(wclData.reportCastEvents.data as Array<wclCastEvent>, event => String(event.abilityGameID));
		professionSpells.forEach(spell => {
			const spellEvents: Array<wclCastEvent> = castEventsBySpellId[spell.id] || [];
			spellEvents.forEach(event => {
				const sourcePlayer = wclPlayers.find(player => player.id == event.sourceID);
				if (sourcePlayer && !sourcePlayer.inferredProfessions.includes(spell.profession)) {
					console.log(`Inferring player ${sourcePlayer.name} has profession ${professionNames.get(spell.profession)} from ${spell.name} event`);
					sourcePlayer.inferredProfessions.push(spell.profession);
				}
			});
		});

		wclPlayers.forEach(player => {
			let professions = distinct(player.inferredProfessions.concat(player.player.getGear().getProfessionRequirements()));
			if (professions.length == 0) {
				professions = [Profession.Engineering, Profession.Jewelcrafting];
			} else if (professions.length == 1) {
				if (professions[0] != Profession.Engineering) {
					professions.push(Profession.Engineering);
				} else {
					professions.push(Profession.Jewelcrafting);
				}
			}
			player.player.setProfessions(eventID, professions);
		});
	}

	private inferAssignments(eventID: EventID, wclData: any, wclPlayers: WCLSimPlayer[]) {
		const castEventsBySpellId = bucket(wclData.fightCastEvents.data as Array<wclCastEvent>, event => String(event.abilityGameID));
		externalCDSpells.forEach(spell => {
			const spellEvents: Array<wclCastEvent> = castEventsBySpellId[spell.id] || [];
			spellEvents.forEach(event => {
				const sourcePlayer = wclPlayers.find(player => player.id == event.sourceID);
				const targetPlayer = wclPlayers.find(player => player.id == event.targetID);
				if (sourcePlayer && targetPlayer && sourcePlayer.player.getClass() == spell.class) {
					const specOptions = spell.applyFunc(sourcePlayer.player, targetPlayer.toUnitReference());
					sourcePlayer.player.setSpecOptions(eventID, specOptions);
					console.log(`Inferring player ${sourcePlayer.name} is targeting ${targetPlayer.name} with ${spell.name} from cast event`);
				}
			});
		});
	}

	// Assigns the raidIndex field for all players.
	private inferPartyComposition(eventID: EventID, wclData: any, wclPlayers: WCLSimPlayer[]) {
		const setPlayersInParty = (player1: WCLSimPlayer, player2: WCLSimPlayer, reason: string) => {
			if (player1.addPlayerInParty(player2) || player2.addPlayerInParty(player1)) {
				console.log(`Inferring players ${player1.name} and ${player2.name} in same party from ${reason} event`);
			}
		};

		const healEventsBySpellId = bucket(wclData.fightHealEvents.data as Array<wclHealEvent>, event => String(event.abilityGameID));

		// These spells only affect players in the same party as the caster.
		samePartyHealingSpells.forEach(spell => {
			const spellEvents: Array<wclHealEvent> = healEventsBySpellId[spell.id] || [];
			spellEvents.forEach(event => {
				const sourcePlayer = wclPlayers.find(player => player.id == event.sourceID);
				const targetPlayer = wclPlayers.find(player => player.id == event.targetID);
				if (sourcePlayer && targetPlayer) {
					setPlayersInParty(sourcePlayer, targetPlayer, spell.name);
				}
			});
		});

		// Prayer of Healing is a bit different, we can infer that players who are targeted at the same time are in a group.
		otherPartyHealingSpells.forEach(spell => {
			const spellEvents: Array<wclHealEvent> = healEventsBySpellId[spell.id] || [];
			const spellEventsByTimestamp = bucket(spellEvents, event => String(event.timestamp) + String(event.sourceID));
			for (const [_timestamp, eventsAtTime] of Object.entries(spellEventsByTimestamp)) {
				const spellTargets = eventsAtTime.map(event => wclPlayers.find(player => player.id == event.targetID));
				for (let i = 0; i < spellTargets.length; i++) {
					for (let j = 0; j < spellTargets.length; j++) {
						if (i != j && spellTargets[i] && spellTargets[j]) {
							setPlayersInParty(spellTargets[i]!, spellTargets[j]!, spell.name);
						}
					}
				}
			}
		});

		wclData.combatantInfoEvents.data.forEach((combatantInfo: wclCombatantInfoEvent) => {
			const targetPlayer = wclPlayers.find(player => player.id == combatantInfo.sourceID);
			combatantInfo.auras
				.filter(aura => aura.ability == 28878)
				.forEach(aura => {
					const sourcePlayer = wclPlayers.find(player => player.id == aura.source);
					if (sourcePlayer && targetPlayer) {
						setPlayersInParty(sourcePlayer, targetPlayer, 'Heroic Presence');
					}
				});
		});

		// Assign players with same-group inferences.
		let inferredPlayers = wclPlayers.filter(player => player.playersInParty.length > 0);
		let nextEmptyPartyIdx = 0;
		while (inferredPlayers.length > 0) {
			// Find all the players in the same party as the first player.
			let partyMembers = [inferredPlayers[0]].concat(inferredPlayers[0].playersInParty);
			let numMembers = 0;
			while (partyMembers.length != numMembers) {
				numMembers = partyMembers.length;
				partyMembers = distinct(partyMembers.map(member => [member].concat(member.playersInParty)).flat());
			}

			// Assign these members to an empty party.
			const partyIdx = nextEmptyPartyIdx;
			nextEmptyPartyIdx++;
			partyMembers.forEach((member, i) => {
				member.raidIndex = partyIdx * 5 + i;
			});

			inferredPlayers = inferredPlayers.filter(player => !partyMembers.includes(player));
		}

		// Assign remaining players into open slots.
		const allRaidIndexes = [...Array(40).keys()];
		const nextFreeIndex = () => allRaidIndexes.find(idx => !wclPlayers.some(p => p.raidIndex == idx)) ?? -1;
		wclPlayers
			.filter(player => player.raidIndex == -1)
			.forEach(player => {
				const nextIdx = nextFreeIndex();
				if (nextIdx == -1) {
					throw new Error('Invalid next idx');
				}
				player.raidIndex = nextIdx;
			});
	}

	private getEncounterProto(wclData: any): EncounterProto {
		const fight: { startTime: number, endTime: number, id: number, name: string } = wclData.fights[0];

		const encounter = EncounterProto.create({
			duration: (fight.endTime - fight.startTime) / 1000,
			targets: [],
		});

		// Use the preset encounter if it exists.
		let closestEncounterPreset = this.simUI.sim.db.getAllPresetEncounters().find(enc => enc.path.includes(fight.name));
		if (closestEncounterPreset && closestEncounterPreset.targets.length) {
			closestEncounterPreset.targets
				.map(mob => mob.target as TargetProto)
				.filter(target => target !== undefined)
				.forEach(target => encounter.targets.push(target));
		}

		// Build a manual target list if no preset encounter exists.
		if (encounter.targets.length === 0) {
			encounter.targets.push(Encounter.defaultTargetProto());
		}

		return encounter;
	}

	private getRaidProto(wclPlayers: WCLSimPlayer[]): RaidProto {
		const raid = RaidProto.create({
			parties: [...new Array(MAX_NUM_PARTIES).keys()].map(_party => PartyProto.create({
				players: [...new Array(5).keys()].map(_player => PlayerProto.create()),
			})),
		});

		wclPlayers
			.forEach(player => {
				const positionInParty = player.raidIndex % 5;
				const partyIdx = (player.raidIndex - positionInParty) / 5;
				const playerProto = player.player.toProto();
				raid.parties[partyIdx].players[positionInParty] = playerProto;

				if (isTankSpec(playerToSpec(playerProto))) {
					raid.tanks.push(player.toUnitReference());
				}
			});

		return raid;
	}
}

class WCLSimPlayer {
	public readonly data: wclPlayer;
	public readonly id: number;
	public readonly name: string;
	public readonly type: string;
	public raidIndex: number = -1;

	private readonly simUI: RaidSimUI;
	private readonly fullType: string;
	private readonly spec: Spec | null;

	readonly player: Player<any>;
	readonly preset: RaidSimPreset<any>;

	inferredProfessions: Array<Profession> = [];

	readonly playersInParty: Array<WCLSimPlayer> = [];

	constructor(data: wclPlayer, simUI: RaidSimUI, eventID: EventID) {
		this.simUI = simUI;
		this.data = data;

		this.name = data.name;
		this.id = data.id;
		this.type = data.type;

		const wclSpec = data.icon.split('-')[1];
		this.fullType = this.type + wclSpec;
		console.log(`WCL spec: ${this.fullType}`);

		const foundSpec = fullTypeToSpec[this.fullType] ?? null;
		if (foundSpec == null) {
			throw new Error('Player type not implemented: ' + this.fullType);
		}
		this.spec = foundSpec;
		this.player = new Player(this.spec, simUI.sim);

		this.preset = WCLSimPlayer.getMatchingPreset(foundSpec, data.talents);
		if (this.preset === undefined) {
			throw new Error('Could not find matching preset: ' + JSON.stringify({
				'name': this.name,
				'type': this.fullType,
				'talents': data.talents,
			}).toString());
		}

		// Apply preset defaults.
		this.player.applySharedDefaults(eventID);
		this.player.setTalentsString(eventID, this.preset.talents.talentsString);
		this.player.setGlyphs(eventID, this.preset.talents.glyphs!);
		this.player.setConsumes(eventID, this.preset.consumes);
		this.player.setSpecOptions(eventID, this.preset.specOptions);
		this.player.setProfessions(eventID, [Profession.Engineering, Profession.Jewelcrafting]);

		// Apply settings from report data.
		this.player.setName(eventID, data.name);
		this.player.setGear(eventID, simUI.sim.db.lookupEquipmentSpec(EquipmentSpec.create({
			items: data.gear.map(gear => ItemSpec.create({
				id: gear.id,
				enchant: gear.permanentEnchant,
				gems: gear.gems ? gear.gems.map(gemInfo => gemInfo.id) : [],
			})),
		})));
	}

	private static getMatchingPreset(spec: Spec, talents: wclTalents[]): RaidSimPreset<Spec> {
		const matchingPresets = playerPresets.filter((preset) => preset.spec == spec);
		let presetIdx = 0;

		if (matchingPresets && matchingPresets.length > 1) {
			let distance = 999;
			// Search talents and find the preset that the players talents most closely match.
			matchingPresets.forEach((preset, i) => {
				const presetTalents = getTalentTreePoints(preset.talents.talentsString);
				// Diff the distance to the preset.
				const newDistance = presetTalents.reduce((acc, v, i) => acc += Math.abs(talents[i]?.guid - presetTalents[i]), 0);

				// If this is the best distance, assign this preset.
				if (newDistance < distance) {
					presetIdx = i;
					distance = newDistance;
				}
			});
		}
		return matchingPresets[presetIdx];
	}

	public toUnitReference(): UnitReference {
		return UnitReference.create({
			type: UnitReference_Type.Player,
			index: this.raidIndex,
		});
	}

	public addPlayerInParty(other: WCLSimPlayer): boolean {
		if (other != this && !this.playersInParty.includes(other)) {
			this.playersInParty.push(other);
			return true;
		}
		return false;
	}
}

const fullTypeToSpec: Record<string, Spec> = {
	'DeathKnightBlood': Spec.SpecTankDeathknight,
	'DeathKnightLichborne': Spec.SpecTankDeathknight,
	'DeathKnightRuneblade': Spec.SpecDeathknight,
	'DeathKnightBloodDPS': Spec.SpecDeathknight,
	'DeathKnightFrost': Spec.SpecDeathknight,
	'DeathKnightUnholy': Spec.SpecDeathknight,
	'DruidBalance': Spec.SpecBalanceDruid,
	'DruidFeral': Spec.SpecFeralDruid,
	'DruidWarden': Spec.SpecFeralTankDruid,
	'DruidGuardian': Spec.SpecFeralTankDruid,
	'DruidRestoration': Spec.SpecRestorationDruid,
	'HunterBeastMastery': Spec.SpecHunter,
	'HunterSurvival': Spec.SpecHunter,
	'HunterMarksmanship': Spec.SpecHunter,
	'MageArcane': Spec.SpecMage,
	'MageFire': Spec.SpecMage,
	'MageFrost': Spec.SpecMage,
	'PaladinHoly': Spec.SpecHolyPaladin,
	'PaladinJusticar': Spec.SpecProtectionPaladin,
	'PaladinProtection': Spec.SpecProtectionPaladin,
	'PaladinRetribution': Spec.SpecRetributionPaladin,
	'PriestHoly': Spec.SpecHealingPriest,
	'PriestDiscipline': Spec.SpecHealingPriest,
	'PriestShadow': Spec.SpecShadowPriest,
	'PriestSmite': Spec.SpecSmitePriest,
	'RogueAssassination': Spec.SpecRogue,
	'RogueCombat': Spec.SpecRogue,
	'RogueSubtlety': Spec.SpecRogue,
	'ShamanElemental': Spec.SpecElementalShaman,
	'ShamanEnhancement': Spec.SpecEnhancementShaman,
	'ShamanRestoration': Spec.SpecRestorationShaman,
	'WarlockDestruction': Spec.SpecWarlock,
	'WarlockAffliction': Spec.SpecWarlock,
	'WarlockDemonology': Spec.SpecWarlock,
	'WarriorArms': Spec.SpecWarrior,
	'WarriorFury': Spec.SpecWarrior,
	'WarriorChampion': Spec.SpecWarrior,
	'WarriorWarrior': Spec.SpecWarrior,
	'WarriorGladiator': Spec.SpecWarrior,
	'WarriorProtection': Spec.SpecProtectionWarrior,
};

// Spells which imply a specific Race.
const racialSpells: Array<{ id: number, name: string, race: Race }> = [
	{ id: 25046, name: 'Arcane Torrent (Energy)', race: Race.RaceBloodElf },
	{ id: 28730, name: 'Arcane Torrent (Mana)', race: Race.RaceBloodElf },
	{ id: 50613, name: 'Arcane Torrent (Runic Power)', race: Race.RaceBloodElf },
	{ id: 26297, name: 'Berserking', race: Race.RaceTroll },
	{ id: 20572, name: 'Blood Fury (AP)', race: Race.RaceOrc },
	{ id: 33697, name: 'Blood Fury (AP+SP)', race: Race.RaceOrc },
	{ id: 33702, name: 'Blood Fury (SP)', race: Race.RaceOrc },
	{ id: 20589, name: 'Escape Artist', race: Race.RaceGnome },
	{ id: 20594, name: 'Stoneform', race: Race.RaceDwarf },
	{ id: 20549, name: 'War Stomp', race: Race.RaceTauren },
	{ id: 7744, name: 'Will of the Forsaken', race: Race.RaceUndead },
	{ id: 59752, name: 'Will to Survive', race: Race.RaceHuman },
];

// Spells which imply a specific Profession.
const professionSpells: Array<{ id: number, name: string, profession: Profession }> = [
	{ id: 55503, name: 'Lifeblood', profession: Profession.Herbalism },
	{ id: 50305, name: 'Skinning', profession: Profession.Skinning },
];

const externalCDSpells: Array<{ id: number, name: string, class: Class, applyFunc: (player: Player<any>, raidTarget: UnitReference) => SpecOptions<any> }> = [
	{
		id: 29166, name: 'Innervate', class: Class.ClassDruid, applyFunc: (player: Player<any>, raidTarget: UnitReference) => {
			const options = player.getSpecOptions() as SpecOptions<DruidSpecs>;
			options.innervateTarget = raidTarget;
			return options;
		}
	},
	{
		id: 10060, name: 'Power Infusion', class: Class.ClassPriest, applyFunc: (player: Player<any>, raidTarget: UnitReference) => {
			const options = player.getSpecOptions() as SpecOptions<PriestSpecs>;
			options.powerInfusionTarget = raidTarget;
			return options;
		}
	},
	{
		id: 57933, name: 'Tricks of the Trade', class: Class.ClassRogue, applyFunc: (player: Player<any>, raidTarget: UnitReference) => {
			const options = player.getSpecOptions() as SpecOptions<RogueSpecs>;
			options.tricksOfTheTradeTarget = raidTarget;
			return options;
		}
	},
	{
		id: 49016, name: 'Unholy Frenzy', class: Class.ClassDeathknight, applyFunc: (player: Player<any>, raidTarget: UnitReference) => {
			const options = player.getSpecOptions() as SpecOptions<DeathknightSpecs>;
			options.unholyFrenzyTarget = raidTarget;
			return options;
		}
	},
];

// Healing spells which only affect the caster's party.
const samePartyHealingSpells: Array<{ id: number, name: string }> = [
	{ id: 52042, name: 'Healing Stream Totem' },
	{ id: 48076, name: 'Holy Nova' },
	{ id: 48445, name: 'Tranquility' },
	{ id: 15290, name: 'Vampiric Embrace' },
];

// Healing spells which only affect a single party, but not necessarily the caster's party.
const otherPartyHealingSpells: Array<{ id: number, name: string }> = [
	{ id: 48072, name: 'Prayer of Healing' },
];

interface wclUrlData {
	reportID: string,
	fightID: string,
}

interface wclCastEvent {
	type: 'cast',
	timestamp: number;
	sourceID: number;
	targetID: number;
	abilityGameID: number;
	fight: number;
}

interface wclHealEvent {
	type: 'heal',
	timestamp: number;
	sourceID: number;
	targetID: number;
	abilityGameID: number;
	fight: number;
	amount: number;
}

interface wclCombatantInfoEvent {
	type: 'combatantinfo';
	sourceID: number;
	auras: {
		source: number;
		ability: number;
		name: string;
	}[];
}

interface wclRateLimitData {
	limitPerHour: number,
	pointsSpentThisHour: number,
	pointsResetIn: number
}

// Typed interface for WCL talents
interface wclTalents {
	name: string;
	guid: number;
	type: number;
	abilityIcon: string;
}

// Typed interface for WCL Gems
interface wclGems {
	id: number;
	itemLevel: number;
	icon: string;
}

// Typed interface for WCL Gear
interface wclGear {
	id: number;
	slot: number;
	quality: number;
	icon: string;
	name: string;
	itemLevel: number;
	permanentEnchant: number;
	permanentEnchantName: string;
	temporaryEnchant: number;
	gems?: wclGems[];
}

// Typed interface for WCL Player Data
interface wclPlayer {
	name: string;
	id: number;
	guid?: number;
	type: string; // Paladin, Mage, etc.
	icon: string; // Paladin-Justicar, Mage-Fire, etc.
	itemLevel?: number;
	total?: number;
	activeTime?: number;
	activeTimeReduced?: number;
	abilities?: unknown[]; // Don't care about abilities.
	damageAbilities?: unknown[];
	targets?: unknown[];
	talents: wclTalents[];
	gear: wclGear[];
}

interface _wclAura {
	name: string;
	id: number;
	guid: number;
	type: string;
	icon: string;
	totalUptime: number;
	totalUses: number;
	bands: {
		startTime: number,
		endTime: number,
	}[];
}
