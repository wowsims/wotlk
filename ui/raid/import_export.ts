import { Exporter } from '../core/components/exporters';
import { Importer } from '../core/components/importers';
import { MAX_PARTY_SIZE } from '../core/party';
import { RaidSimSettings } from '../core/proto/ui';
import { EventID, TypedEvent } from '../core/typed_event';
import { Party as PartyProto, Player as PlayerProto, Raid as RaidProto } from '../core/proto/api';
import {
	Class,
	Encounter as EncounterProto,
	EquipmentSpec,
	Faction,
	ItemSpec,
	MobType,
	Profession,
	Race,
	RaidTarget,
	Spec,
	Target as TargetProto,
} from '../core/proto/common';
import { nameToClass } from '../core/proto_utils/names';
import {
	DruidSpecs,
	DeathknightSpecs,
	MageSpecs,
	PriestSpecs,
	RogueSpecs,
	getTalentTreePoints,
	makeDefaultBlessings,
	specTypeFunctions,
	withSpecProto,
	isTankSpec,
	playerToSpec,
} from '../core/proto_utils/utils';
import { MAX_NUM_PARTIES } from '../core/raid';
import { Player } from '../core/player';
import { Target } from '../core/target';
import { bucket, distinct, sortByProperty } from '../core/utils';

import { playerPresets, PresetSpecSettings } from './presets';
import { RaidSimUI } from './raid_sim_ui';

declare var $: any;
declare var tippy: any;

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
		super(parent, simUI, 'JSON Export', true);
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
		const match = url.match(/classic\.warcraftlogs\.com\/reports\/([a-zA-Z0-9:]+)(#.*fight=((\d+)|(last)))?/);
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
		const rateLimit = await this.getRateLimit();

		// Schema for WCL API here: https://www.warcraftlogs.com/v2-api-docs/warcraft/
		const reportDataQuery = `{
			reportData {
				report(code: "${urlData.reportID}") {
					guild {
						name faction {id}
					}
					playerDetails: table(fightIDs: [${urlData.fightID}], endTime: 99999999, dataType: Casts, killType: All, viewBy: Default)
					allEvents: events(fightIDs: [${urlData.fightID}], dataType:All, endTime: 99999999, abilityID: 54172, limit: 50) { data }
					fights(fightIDs: [${urlData.fightID}]) {
						startTime, endTime, id, name
					}

					arcaneTorrentEnergy: table(dataType:Casts, endTime: 99999999, abilityID: 25046)
					arcaneTorrentMana: table(dataType:Casts, endTime: 99999999, abilityID: 28730)
					arcaneTorrentRunicPower: table(dataType:Casts, endTime: 99999999, abilityID: 50613)
					berserking: table(dataType:Casts, endTime: 99999999, abilityID: 26297)
					bloodFuryAp: table(dataType:Casts, endTime: 99999999, abilityID: 20572)
					bloodFuryApSp: table(dataType:Casts, endTime: 99999999, abilityID: 33697)
					bloodFurySp: table(dataType:Casts, endTime: 99999999, abilityID: 33702)
					escapeArtist: table(dataType:Casts, endTime: 99999999, abilityID: 20589)
					stoneform: table(dataType:Casts, endTime: 99999999, abilityID: 20594)
					warStomp: table(dataType:Casts, endTime: 99999999, abilityID: 20549)
					willOfTheForsaken: table(dataType:Casts, endTime: 99999999, abilityID: 7744)
					willToSurvive: table(dataType:Casts, endTime: 99999999, abilityID: 59752)

					lifeblood: table(dataType:Casts, endTime: 99999999, abilityID: 55503)
					skinning: table(dataType:Casts, endTime: 99999999, abilityID: 50305)

					innervates: table(fightIDs: [${urlData.fightID}], dataType:Casts, endTime: 99999999, sourceClass: "Druid", abilityID: 29166),
					powerInfusion: table(fightIDs: [${urlData.fightID}], dataType:Casts, endTime: 99999999, sourceClass: "Priest", abilityID: 10060)
					tricksOfTheTrade: table(fightIDs: [${urlData.fightID}], dataType:Casts, endTime: 99999999, sourceClass: "Rogue", abilityID: 57933)
					unholyFrenzy: table(fightIDs: [${urlData.fightID}], dataType:Casts, endTime: 99999999, sourceClass: "DeathKnight", abilityID: 49016)

					divineStorm: events(fightIDs: [${urlData.fightID}], dataType:Healing, endTime: 99999999, sourceClass: "Paladin", abilityID: 54172, limit: 50) { data }
					healingStreamTotem: events(fightIDs: [${urlData.fightID}], dataType:Healing, endTime: 99999999, sourceClass: "Shaman", abilityID: 52042, limit: 50) { data }
					holyNova: events(fightIDs: [${urlData.fightID}], dataType:Healing, endTime: 99999999, sourceClass: "Priest", abilityID: 48076, limit: 100) { data }
					prayerOfHealing: events(fightIDs: [${urlData.fightID}], dataType:Healing, endTime: 99999999, sourceClass: "Priest", abilityID: 48072) { data }
					tranquility: events(fightIDs: [${urlData.fightID}], dataType:Healing, endTime: 99999999, sourceClass: "Druid", abilityID: 48445, limit: 50) { data }
					vampiricEmbrace: events(fightIDs: [${urlData.fightID}], dataType:Healing, endTime: 99999999, sourceClass: "Priest", abilityID: 15290, limit: 50) { data }
				}
			}
		}`;
		const reportData = await this.queryWCL(reportDataQuery);

		// Process the report data.
		const wclData = reportData.data.reportData.report; // TODO: Typings?
		const playerData: wclPlayer[] = wclData.playerDetails.data.entries;

		// If defined in log, use that faction. Otherwise default to UI setting.
		const faction = (wclData.guild?.faction?.id || this.simUI.raidPicker?.getCurrentFaction() || Faction.Horde) as Faction;

		TypedEvent.freezeAllAndDo(() => {
			const eventID = TypedEvent.nextEventID();
			const wclPlayers = playerData.map(wclPlayer => new WCLSimPlayer(wclPlayer, this.simUI, faction, eventID));
			this.inferRace(eventID, wclData, wclPlayers);
			this.inferProfessions(eventID, wclData, wclPlayers);
			this.inferAssignments(eventID, wclData, wclPlayers);
			this.inferPartyComposition(eventID, wclData, wclPlayers);
			const numPaladins = playerData.filter(player => player.type == 'Paladin').length;
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
		wclData.allEvents.data.filter((ev: any) => ev.type == 'combatantinfo').forEach((combatantInfo: wclCombatantInfoEventData) => {
			const targetPlayer = wclPlayers.find(player => player.id == combatantInfo.sourceID);
			combatantInfo.auras
				.filter(aura => aura.ability == 28878)
				.forEach(aura => {
					const sourcePlayer = wclPlayers.find(player => player.id == aura.source);
					if (sourcePlayer && targetPlayer) {
						sourcePlayer.player.setRace(eventID, Race.RaceDraenei);
					}
				});
		});

		this.parseCastData(wclPlayers, wclData.arcaneTorrentMana.data.entries).forEach(cast => {
			cast.player.player.setRace(eventID, Race.RaceBloodElf);
		});
		this.parseCastData(wclPlayers, wclData.arcaneTorrentEnergy.data.entries).forEach(cast => {
			cast.player.player.setRace(eventID, Race.RaceBloodElf);
		});
		this.parseCastData(wclPlayers, wclData.arcaneTorrentRunicPower.data.entries).forEach(cast => {
			cast.player.player.setRace(eventID, Race.RaceBloodElf);
		});
		this.parseCastData(wclPlayers, wclData.berserking.data.entries).forEach(cast => {
			cast.player.player.setRace(eventID, Race.RaceTroll);
		});
		this.parseCastData(wclPlayers, wclData.bloodFuryAp.data.entries).forEach(cast => {
			cast.player.player.setRace(eventID, Race.RaceOrc);
		});
		this.parseCastData(wclPlayers, wclData.bloodFuryApSp.data.entries).forEach(cast => {
			cast.player.player.setRace(eventID, Race.RaceOrc);
		});
		this.parseCastData(wclPlayers, wclData.bloodFurySp.data.entries).forEach(cast => {
			cast.player.player.setRace(eventID, Race.RaceOrc);
		});
		this.parseCastData(wclPlayers, wclData.escapeArtist.data.entries).forEach(cast => {
			cast.player.player.setRace(eventID, Race.RaceGnome);
		});
		this.parseCastData(wclPlayers, wclData.stoneform.data.entries).forEach(cast => {
			cast.player.player.setRace(eventID, Race.RaceDwarf);
		});
		this.parseCastData(wclPlayers, wclData.warStomp.data.entries).forEach(cast => {
			cast.player.player.setRace(eventID, Race.RaceTauren);
		});
		this.parseCastData(wclPlayers, wclData.willOfTheForsaken.data.entries).forEach(cast => {
			cast.player.player.setRace(eventID, Race.RaceUndead);
		});
		this.parseCastData(wclPlayers, wclData.willToSurvive.data.entries).forEach(cast => {
			cast.player.player.setRace(eventID, Race.RaceHuman);
		});
	}

	private inferProfessions(eventID: EventID, wclData: any, wclPlayers: WCLSimPlayer[]) {
		this.parseCastData(wclPlayers, wclData.lifeblood.data.entries).forEach(cast => {
			cast.player.inferredProfessions.push(Profession.Herbalism);
		});
		this.parseCastData(wclPlayers, wclData.skinning.data.entries).forEach(cast => {
			cast.player.inferredProfessions.push(Profession.Skinning);
		});

		wclPlayers.forEach(player => {
			let professions = player.inferredProfessions.concat(player.player.getGear().getProfessionRequirements());
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
		this.parseCastData(wclPlayers, wclData.innervates.data.entries).forEach(cast => {
			if (cast.target && cast.player.player.getClass() == Class.ClassDruid) {
				const player = cast.player.player as Player<DruidSpecs>;
				const options = player.getSpecOptions();
				options.innervateTarget = cast.target.toRaidTarget();
				player.setSpecOptions(eventID, options);
			}
		});
		this.parseCastData(wclPlayers, wclData.powerInfusion.data.entries).forEach(cast => {
			if (cast.target && cast.player.player.getClass() == Class.ClassPriest) {
				const player = cast.player.player as Player<PriestSpecs>;
				const options = player.getSpecOptions();
				options.powerInfusionTarget = cast.target.toRaidTarget();
				player.setSpecOptions(eventID, options);
			}
		});
		this.parseCastData(wclPlayers, wclData.tricksOfTheTrade.data.entries).forEach(cast => {
			if (cast.target && cast.player.player.getClass() == Class.ClassRogue) {
				const player = cast.player.player as Player<RogueSpecs>;
				const options = player.getSpecOptions();
				options.tricksOfTheTradeTarget = cast.target.toRaidTarget();
				player.setSpecOptions(eventID, options);
			}
		});
		this.parseCastData(wclPlayers, wclData.unholyFrenzy.data.entries).forEach(cast => {
			if (cast.target && cast.player.player.getClass() == Class.ClassDeathknight) {
				const player = cast.player.player as Player<DeathknightSpecs>;
				const options = player.getSpecOptions();
				options.unholyFrenzyTarget = cast.target.toRaidTarget();
				player.setSpecOptions(eventID, options);
			}
		});
	}

	private parseCastData(wclPlayers: WCLSimPlayer[], castData: wclBuffCastsData[]): { player: WCLSimPlayer, target: WCLSimPlayer|null }[] {
		const playerCasts: { player: WCLSimPlayer, target: WCLSimPlayer|null }[] = [];
		if (castData.length) {
			castData.forEach(cast => {
				const sourcePlayer = wclPlayers.find((player) => player.name == cast.name);
				const targetPlayer = wclPlayers.find((player) => player.name == cast.targets[0].name) || null;

				if (sourcePlayer) {
					playerCasts.push({ player: sourcePlayer, target: targetPlayer });
				}
			});
		}
		return playerCasts;
	}

	// Assigns the raidIndex field for all players.
	private inferPartyComposition(eventID: EventID, wclData: any, wclPlayers: WCLSimPlayer[]) {
		const parseHealEventData = (healEventData: wclHealEventData[]) => {
			return healEventData.map(event => {
				const sourcePlayer = wclPlayers.find(player => player.id == event.sourceID);
				const targetPlayer = wclPlayers.find(player => player.id == event.targetID);
				if (sourcePlayer && targetPlayer) {
					return { source: sourcePlayer, target: targetPlayer, timestamp: event.timestamp };
				} else {
					return null;
				}
			}).filter(data => data != null) as Array<{source: WCLSimPlayer, target: WCLSimPlayer, timestamp: number}>;
		};

		const inferPartyFromGroupHealData = (healEventData: wclHealEventData[], spellName: string) => {
			parseHealEventData(healEventData).forEach(event => {
				if (event.source.addPlayerInParty(event.target) || event.target.addPlayerInParty(event.source)) {
					console.log(`Inferring players ${event.source.name} and ${event.target.name} in same party from ${spellName} healing event`);
				}
			});
		};
		inferPartyFromGroupHealData(wclData.divineStorm.data, 'Divine Storm');
		inferPartyFromGroupHealData(wclData.healingStreamTotem.data, 'Healing Stream Totem');
		inferPartyFromGroupHealData(wclData.tranquility.data, 'Tranquility');
		inferPartyFromGroupHealData(wclData.vampiricEmbrace.data, 'Vampiric Embrace');

		// Prayer of Healing is a bit different, we can infer that players who are targeted at the same time are in a group.
		const pohByTimestamp = bucket(parseHealEventData(wclData.prayerOfHealing.data), event => String(event.timestamp));
		for (const [k, pohEvents] of Object.entries(pohByTimestamp)) {
			const pohTargets = pohEvents.map(ev => ev.target);
			for (let i = 0; i < pohTargets.length; i++) {
				for (let j = 0; j < pohTargets.length; j++) {
					if (i != j && pohTargets[i].addPlayerInParty(pohTargets[j]) || pohTargets[j].addPlayerInParty(pohTargets[i])) {
						console.log(`Inferring players ${pohTargets[i].name} and ${pohTargets[j].name} in same party from Prayer of Healing event`);
					}
				}
			}
		}

		wclData.allEvents.data.filter((ev: any) => ev.type == 'combatantinfo').forEach((combatantInfo: wclCombatantInfoEventData) => {
			const targetPlayer = wclPlayers.find(player => player.id == combatantInfo.sourceID);
			combatantInfo.auras
				.filter(aura => aura.ability == 28878)
				.forEach(aura => {
					const sourcePlayer = wclPlayers.find(player => player.id == aura.source);
					if (sourcePlayer && targetPlayer) {
						sourcePlayer.addPlayerInParty(targetPlayer);
						targetPlayer.addPlayerInParty(sourcePlayer);
						console.log(`Inferring players ${sourcePlayer.name} and ${targetPlayer.name} in same party from Heroic Presence`);
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
			encounter.targets.push(Target.defaultProto());
		}

		return encounter;
	}

	private getRaidProto(wclPlayers: WCLSimPlayer[]): RaidProto {
		const raid = RaidProto.create({
			parties: [...new Array(MAX_NUM_PARTIES).keys()].map(p => PartyProto.create({
				players: [...new Array(5).keys()].map(p => PlayerProto.create()),
			})),
		});

		wclPlayers
			.forEach(player => {
				const positionInParty = player.raidIndex % 5;
				const partyIdx = (player.raidIndex - positionInParty) / 5;
				const playerProto = player.player.toProto();
				raid.parties[partyIdx].players[positionInParty] = playerProto;

				if (isTankSpec(playerToSpec(playerProto))) {
					raid.tanks.push(player.toRaidTarget());
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
	private readonly spec: Spec|null;
	private readonly faction: Faction;

	readonly player: Player<any>;
	readonly preset: PresetSpecSettings<any>;

	inferredProfessions: Array<Profession> = [];

	readonly playersInParty: Array<WCLSimPlayer> = [];
	readonly playersNotInParty: Array<WCLSimPlayer> = [];

	constructor(data: wclPlayer, simUI: RaidSimUI, faction: Faction = Faction.Unknown, eventID: EventID) {
		this.simUI = simUI;
		this.data = data;

		this.name = data.name;
		this.id = data.id;
		this.type = data.type;
		this.faction = faction;

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
		this.player.setRace(eventID, this.preset.defaultFactionRaces[this.faction]);
		this.player.setTalentsString(eventID, this.preset.talents.talentsString);
		this.player.setGlyphs(eventID, this.preset.talents.glyphs!);
		this.player.setConsumes(eventID, this.preset.consumes);
		this.player.setRotation(eventID, this.preset.rotation);
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

	private static getMatchingPreset(spec: Spec, talents: wclTalents[]): PresetSpecSettings<Spec> {
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

	public toRaidTarget(): RaidTarget {
		return RaidTarget.create({
			targetIndex: this.raidIndex,
		});
	}

	public addPlayerInParty(other: WCLSimPlayer): boolean {
		if (!this.playersInParty.includes(other)) {
			this.playersInParty.push(other);
			return true;
		}
		return false;
	}

	public addPlayerNotInParty(other: WCLSimPlayer): boolean {
		if (!this.playersNotInParty.includes(other)) {
			this.playersNotInParty.push(other);
			return true;
		}
		return false;
	}
}

const fullTypeToSpec: Record<string, Spec> = {
	'DeathKnightBlood': Spec.SpecTankDeathknight,
	'DeathKnightLichborne': Spec.SpecTankDeathknight,
	'DeathKnightRuneblade': Spec.SpecDeathknight,
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

interface wclUrlData {
	reportID: string,
	fightID: string,
}

interface wclBuffCastsData {
	name: string;
	targets: {
		name: string;
		type: string;
	}[];
}

interface wclHealEventData {
	timestamp: number;
	sourceID: number;
	targetID: number;
	amount: number;
}

interface wclCombatantInfoEventData {
	sourceID: number;
	type: string;
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

interface wclAura {
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
