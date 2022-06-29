import { Exporter } from '/tbc/core/components/exporters.js';
import { Importer } from '/tbc/core/components/importers.js';
import { MAX_PARTY_SIZE } from '/tbc/core/party.js';
import { BuffBot, RaidSimSettings } from '/tbc/core/proto/ui.js';
import { TypedEvent } from '/tbc/core/typed_event.js';
import { Party as PartyProto, Player as PlayerProto, Raid as RaidProto } from '/tbc/core/proto/api.js';
import { Encounter as EncounterProto, EquipmentSpec, ItemSpec, MobType, Spec, Target as TargetProto, RaidTarget, Faction } from '/tbc/core/proto/common.js';
import { nameToClass } from '/tbc/core/proto_utils/names.js';
import { makeDefaultBlessings, specTypeFunctions, withSpecProto, isTankSpec, playerToSpec } from '/tbc/core/proto_utils/utils.js';
import { MAX_NUM_PARTIES } from '/tbc/core/raid.js';
import { Player } from '/tbc/core/player.js';
import { Target } from '/tbc/core/target.js';

import { playerPresets, PresetSpecSettings } from './presets.js';
import { RaidSimUI } from './raid_sim_ui.js';

declare var $: any;
declare var tippy: any;

export function newRaidImporters(simUI: RaidSimUI): HTMLElement {
	const importSettings = document.createElement('div');
	importSettings.classList.add('import-settings', 'sim-dropdown-menu');
	importSettings.innerHTML = `
		<span id='importMenuLink' class='dropdown-toggle fas fa-file-import' role='button' data-toggle='dropdown' aria-haspopup='true' arai-expanded='false'></span>
		<div class='dropdown-menu dropdown-menu-right' aria-labelledby='importMenuLink'>
		</div>
	`;
	const linkElem = importSettings.getElementsByClassName('dropdown-toggle')[0] as HTMLElement;
	tippy(linkElem, {
		'content': 'Import',
		'allowHTML': true,
	});

	const menuElem = importSettings.getElementsByClassName('dropdown-menu')[0] as HTMLElement;
	const addMenuItem = (label: string, experimental: boolean, onClick: () => void) => {
		const itemElem = document.createElement('span');
		itemElem.classList.add('dropdown-item');
		itemElem.textContent = label;
		itemElem.addEventListener('click', onClick);
		menuElem.appendChild(itemElem);
	};

	addMenuItem('Json', false, () => new RaidJsonImporter(menuElem, simUI));
	addMenuItem('WCL', true, () => new RaidWCLImporter(menuElem, simUI));

	return importSettings;
}

export function newRaidExporters(simUI: RaidSimUI): HTMLElement {
	const exportSettings = document.createElement('div');
	exportSettings.classList.add('export-settings', 'sim-dropdown-menu');
	exportSettings.innerHTML = `
		<span id='exportMenuLink' class='dropdown-toggle fas fa-file-export' role='button' data-toggle='dropdown' aria-haspopup='true' arai-expanded='false'></span>
		<div class='dropdown-menu dropdown-menu-right' aria-labelledby='exportMenuLink'>
		</div>
	`;
	const linkElem = exportSettings.getElementsByClassName('dropdown-toggle')[0] as HTMLElement;
	tippy(linkElem, {
		'content': 'Export',
		'allowHTML': true,
	});

	const menuElem = exportSettings.getElementsByClassName('dropdown-menu')[0] as HTMLElement;
	const addMenuItem = (label: string, onClick: () => void) => {
		const itemElem = document.createElement('span');
		itemElem.classList.add('dropdown-item');
		itemElem.textContent = label;
		itemElem.addEventListener('click', onClick);
		menuElem.appendChild(itemElem);
	};

	addMenuItem('Json', () => new RaidJsonExporter(menuElem, simUI));

	return exportSettings;
}

class RaidJsonImporter extends Importer {
	private readonly simUI: RaidSimUI;
	constructor(parent: HTMLElement, simUI: RaidSimUI) {
		super(parent, 'JSON Import', true);
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
		const settings = RaidSimSettings.fromJsonString(data);
		this.simUI.fromProto(TypedEvent.nextEventID(), settings);
		this.close();
	}
}

class RaidJsonExporter extends Exporter {
	private readonly simUI: RaidSimUI;

	constructor(parent: HTMLElement, simUI: RaidSimUI) {
		super(parent, 'JSON Export', true);
		this.simUI = simUI;
		this.init();
	}

	getData(): string {
		return JSON.stringify(RaidSimSettings.toJson(this.simUI.toProto()), null, 2);
	}
}

class RaidWCLImporter extends Importer {

	private queryCounter: number = 0;

	private readonly simUI: RaidSimUI;
	constructor(parent: HTMLElement, simUI: RaidSimUI) {
		super(parent, 'WCL Import', false);
		this.simUI = simUI;
		this.descriptionElem.innerHTML = `
			<p>
				Import entire raid from a WCL report.<br>
				Parties are a best guess based on buffs.<br>
				Double check innervate/PI and paladin buffs in the settings after import.<br>
				Does not support fight=last currently (will default any non-numeric fight ID to be 0)<br>
			</p>
			<p>
				To import, paste the WCL report and fight link (https://classic.warcraftlogs.com/reports/REPORTID#fight=FIGHTID).<br>
				Include the fight ID or else first found fight will be used.<br>
			</p>
		`;
	}

	private getWCLBearerToken(): Promise<string> {

		return fetch('https://classic.warcraftlogs.com/oauth/token', {
			'method': 'POST',
			'headers': {
				'Authorization': 'Basic ' + btoa('963d31c8-7efa-4dde-87cf-1b254a8a2f8c:lRJVhujEEnF96xfUoxVHSpnqKN9v8bTqGEjutsO3'),
			},
			body: new URLSearchParams({
				'grant_type': 'client_credentials',
			}),
		}).then(response => response.json())
			.then(res => res.access_token)
			.catch(err => { // TODO: handle error
				console.error(err);
			});
	}

	private queryWCL(query: string, token: string): Promise<any> {
		const headers = {
			'Content-Type': 'application/json',
			'Authorization': `Bearer ${token}`,
			'Accept': 'application/json',
		};

		const queryURL = `https://classic.warcraftlogs.com/api/v2/client?query=${query}`;

		this.queryCounter++;

		// Query WCL
		return fetch(encodeURI(queryURL), {
			'method': 'GET',
			'headers': headers,
		}).then(res => res.json());
	}

	private getURLInfo(url: string): { reportID: string; fightID: string } {

		let urlInfo = { reportID: '', fightID: '0' };

		if (!url.includes('warcraftlogs.com')) {
			console.error('Invalid WCL URL', url, 'must be from warcraftlogs.com');
			return urlInfo;
		}

		let fightIDIndex = url.indexOf('fight=');
		let reportIDIndex = url.indexOf('/reports/');

		if (reportIDIndex === -1) {
			console.error('Could not find report ID in URL', url);
			return urlInfo;
		}

		reportIDIndex += 9; // 9 = length of '/reports/'
		const reportIDLength = 16;

		if (fightIDIndex !== -1) {
			fightIDIndex += 6; // 6 = length of 'fight='

			let fightID = parseInt(url.substring(fightIDIndex), 10);

			if (isNaN(fightID)) {
				fightID = 0;
			}

			urlInfo.fightID = fightID.toString();
		} else {
			console.warn('Could not find fight ID in URL', url, 'defaulting to fight 0');
		}

		urlInfo.reportID = url.substring(reportIDIndex, reportIDIndex + reportIDLength) ?? '';

		return urlInfo
	}

	private getRateLimit(token: string): Promise<wclRateLimitData> {
		const query = `
	  {
	    rateLimitData {
	      limitPerHour, pointsSpentThisHour, pointsResetIn
	    }
	  }`;
		return this.queryWCL(query, token)
			.then((res) => res['data']['rateLimitData'] as wclRateLimitData);
	}

	async onImport(importLink: string) {
		this.importButton.disabled = true;
		this.rootElem.style.cursor = 'wait';
		this.doImport(importLink).then(() => {
			this.importButton.disabled = false
			this.rootElem.style.removeProperty('cursor');
		})
	}

	async doImport(importLink: string) {
		if (!importLink.length) {
			console.error('No import link provided!');
			return;
		}

		let urlInfo = this.getURLInfo(importLink);

		const reportID = urlInfo.reportID;
		const fightID = urlInfo.fightID;

		if (!reportID.length) {
			console.error('Could not find report ID in URL', importLink);
			return;
		}

		// Clear the raid out to avoid any taint issues.
		this.simUI.clearRaid(TypedEvent.nextEventID());

		const token = await this.getWCLBearerToken();

		const rateLimitBuffer = 30; // WCL Query point buffer
		const rateLimitStart: wclRateLimitData = await this.getRateLimit(token);

		// Slower but more accurate way to generate the raid sim.
		// Generates players into the groups that they were in during the fight.
		// If the rate limit is close to max, then it will create the raid parties 'randomly'.
		let experimentalGenerateParties: boolean = rateLimitStart.pointsSpentThisHour + rateLimitBuffer < rateLimitStart.limitPerHour;

		console.info('Importing WCL report', reportID, 'fight', fightID, 'Generate Parties:', experimentalGenerateParties);

		const reportDataQuery = `
				{
					reportData {
						report(code: "${reportID}") {
							guild {
								name faction {id}
							}
							playerDetails: table(fightIDs: [${fightID}], endTime: 99999999, dataType: Casts, killType: All, viewBy: Default)
							fights(fightIDs: [${fightID}]) {
								startTime, endTime, id, name
							}
							innervates: table(fightIDs: [${fightID}], dataType:Casts, endTime: 99999999, sourceClass: "Druid", abilityID: 29166),
							powerInfusion: table(fightIDs: [${fightID}], dataType:Casts, endTime: 99999999, sourceClass: "Priest", abilityID: 10060)
						}
					}
				}
				`;

		const reportData = await this.queryWCL(reportDataQuery, token);
		if (reportData.errors != undefined && reportData.errors != null && reportData.errors.length) {
			const errorData = reportData.errors.reduce((accumulator: string, error: any) => {
				return accumulator + error.message;
			}, "");
			alert("Failed to import: " + errorData);
			return;
		}
		// Process the report data.
		const wclData = reportData.data.reportData.report; // TODO: Typings?

		const guildData = wclData.guild;
		const playerData: wclPlayer[] = wclData.playerDetails.data.entries;
		const innervateData: wclBuffCastsData[] = wclData.innervates.data.entries;
		const powerInfusionData: wclBuffCastsData[] = wclData.powerInfusion.data.entries;

		// Set up the general variables we need for import to be successful.
		const fight: { startTime: number, endTime: number, id: number, name: string } = wclData.fights[0];
		const startTime: number = fight.startTime;
		const endTime: number = fight.endTime;

		// Default to UI setting
		let faction = this.simUI.raidPicker?.getCurrentFaction();

		// If defined in log, use that faction.
		if (guildData != null) {
			faction = guildData.faction.id as Faction;
		}

		// Fallback if UI is broken and log has no faction.
		if (faction == undefined) {
			faction = Faction.Horde;
		}

		const encounter = EncounterProto.create();
		encounter.duration = (endTime - startTime) / 1000;

		encounter.targets = new Array<TargetProto>();

		let closestEncounterPreset = this.simUI.sim.getAllPresetEncounters().find((enc) => enc.path.includes(fight.name));

		// Use the preset encounter if it exists.
		if (closestEncounterPreset && closestEncounterPreset.targets.length) {
			closestEncounterPreset.targets
				.map((mob) => mob.target as TargetProto)
				.filter((target) => target !== undefined)
				.forEach((target) => encounter.targets.push(target));
		}

		// Build a manual target list if no preset encounter exists.
		if (encounter.targets.length === 0) {
			encounter.targets.push(Target.defaultProto());
		}

		const settings = RaidSimSettings.create();
		settings.encounter = encounter;

		const raid = RaidProto.create();
		raid.parties = new Array<PartyProto>();
		settings.raid = raid;

		const buffBots = new Array<BuffBot>();

		// Raid index of players that received innervates
		const wclIDtoRaidIndex = new Map<number, number>();

		const numPaladins = playerData.filter((player) => player.type === 'Paladin').length;

		// Generate an empty set of 3 dimensional arrays for each party. [ party ][ player or buffBot ][ player ]
		let tempParties: WCLSimPlayer[][] = Array.from({ length: MAX_NUM_PARTIES }).map(() => []);

		// Generate the default 5 raid parties & temp parties.
		tempParties.forEach(() => raid.parties.push(PartyProto.create()));

		// Sorts an objectArray by a property. Returns a new array.
		// Can be called recursively.
		const sortByProperty = (objArray: any[], prop: string) => {
			if (!Array.isArray(objArray)) throw new Error('FIRST ARGUMENT NOT AN ARRAY');
			const clone = objArray.slice(0);
			const direct = arguments.length > 2 ? arguments[2] : 1; //Default to ascending
			const propPath = (prop.constructor === Array) ? prop : prop.split('.');
			clone.sort(function(a, b) {
				for (let p in propPath) {
					if (a[propPath[p]] && b[propPath[p]]) {
						a = a[propPath[p]];
						b = b[propPath[p]];
					}
				}
				// convert numeric strings to integers
				a = a.toString().match(/^\d+$/) ? +a : a;
				b = b.toString().match(/^\d+$/) ? +b : b;
				return ((a < b) ? -1 * direct : ((a > b) ? 1 * direct : 0));
			});
			return clone;
		}

		const mappedPlayers = playerData
			.map((wclPlayer) => new WCLSimPlayer(wclPlayer, this.simUI, faction));

		const processBuffCastData = (buffCastData: wclBuffCastsData[]): { player: WCLSimPlayer, target: string }[] => {
			const playerCasts: { player: WCLSimPlayer, target: string }[] = [];
			if (buffCastData.length) {
				buffCastData.forEach((cast) => {
					const sourcePlayer = mappedPlayers.find((player) => player.name === cast.name);
					const targetPlayer = mappedPlayers.find((player) => player.name === cast.targets[0].name);

					// Buff bots do not get PI/Innervates.
					if (sourcePlayer && targetPlayer && !targetPlayer.isBuffBot) {
						playerCasts.push({ player: sourcePlayer, target: targetPlayer.name });
					}
				});
			}
			return playerCasts;
		}

		processBuffCastData(innervateData).forEach((cast) => cast.player.innervateTarget = cast.target);
		processBuffCastData(powerInfusionData).forEach((cast) => cast.player.powerInfusionTarget = cast.target);

		const wclPlayers: WCLSimPlayer[] = sortByProperty(sortByProperty(mappedPlayers, 'type'), 'sortPriority');

		let raidIndex = 0;

		// Sorts buff bots to the end of the array to prevent overwriting them later on.
		const sortBuffBotsLast = (a: WCLSimPlayer, b: WCLSimPlayer) => a.isBuffBot ? 1 : b.isBuffBot ? 1 : 0;

		// Reusable function to add a player to the raid.parties[raidIndex] array.
		const assignPlayerToParty = (player: WCLSimPlayer, raidParty: PartyProto, missing = false) => {

			if (!player) {
				console.error('Cannot assign player to party because player is undefined!');
				return;
			}

			if (!raidParty) {
				console.error('Cannot assign player to party because party is undefined!');
				return;
			}

			if (raidParty.players.length === MAX_PARTY_SIZE) {
				console.error('Cannot assign player to party because party is full!', player, raidParty.players);
				return;
			}

			if (missing) {
				console.warn(`Could not locate a group for ${player.name}, assigning them to an open group.`);
			}

			let buffBot = player.getBuffBot();
			let simPlayer = player.getPlayer();

			if (!buffBot && !simPlayer) {
				console.error('Cannot assign player to party because player data is undefined!', player);
				return;
			}

			wclIDtoRaidIndex.set(player.id, raidIndex);

			if (buffBot) {
				buffBot.raidIndex = raidIndex;
				buffBots.push(buffBot);
				raidParty.players.push(PlayerProto.create());
			} else if (simPlayer) {
				raidParty.players.push(simPlayer);
				if (isTankSpec(playerToSpec(simPlayer))) {
					let rt = RaidTarget.create();
					rt.targetIndex = wclIDtoRaidIndex.get(player.id)!;
					settings.raid!.tanks.push(rt);
				}
			}

			// Just in case this did not get set previously.
			player.partyAssigned = true;

			raidIndex++;
		}

		// if experimental_generate_parties is true, we will generate parties based on the party buffers
		if (experimentalGenerateParties) {

			// We only care about the players who can provide party buffs on logs.
			const partyBuffPlayers = wclPlayers.filter((player) => player.isPartyBuffer);

			// Can't be a forEach because we need to wait for the query to finish on each iteration later on.
			for (const player of partyBuffPlayers) {

				const partyFull = player.partyMembers.length >= MAX_PARTY_SIZE;

				// Skip players that have already been assigned to a party.
				// player.partyAssigned || player.partyFound || player.partyMembers.length > 0
				if (partyFull) {
					continue;
				}

				const auraIDs: number[] = player.getPartyAuraIds();

				if (!auraIDs.length) {
					console.warn('No party aura ids found for partyBuff player ' + player.name);
					continue;
				}

				let auraBuffQueries = auraIDs.map((auraID) => `
				{
					reportData {
						report(code: "${reportID}") {
					table(startTime: ${startTime}, endTime: ${endTime}, sourceID: ${player.id}, abilityID: ${auraID}, fightIDs: [${fightID}],dataType:Buffs,viewBy:Target,hostilityType:Friendlies)
						}
					}
				}`);

				let auraTargets: wclAura[] = [];

				// Can't be a forEach because we need to await each query.
				for (let i = 0; i < auraBuffQueries.length; i++) {

					if (auraTargets.length >= MAX_PARTY_SIZE || partyFull) {
						break;
					}

					let auraQueryRes = await this.queryWCL(auraBuffQueries[i], token);
					if (auraQueryRes) {
						let playerAuras: wclAura[] = auraQueryRes.data?.reportData?.report?.table?.data?.auras ?? [];
						if (playerAuras.length) {

							playerAuras = playerAuras.filter((auraTarget) => auraTarget.type !== 'Pet')
								.sort((a, b) => a.bands[0].startTime - b.bands[0].startTime)
								.filter((auraTarget, index) => index < 5);

							const uniqueAuraTargets = playerAuras.filter((auraTarget) => !auraTargets.some((target) => target.name === auraTarget.name));
							auraTargets.push(...uniqueAuraTargets);
						}
					}
				}

				if (auraTargets.length === 0) {
					continue;
				}

				player.partyFound = true;

				// Only need the member names at this point.
				player.partyMembers = auraTargets.map((auraTarget) => auraTarget.name);

				let partyMembers = wclPlayers
					.filter((raidMember) => player.partyMembers.includes(raidMember.name))
					.filter((raidMember) => !raidMember.partyAssigned);

				const totalPartyMembers = partyMembers.length;

				// Find an empty temp party to assign the members to.
				let partyIndex: number = tempParties.findIndex((party) => party.length < MAX_PARTY_SIZE && party.length < totalPartyMembers);

				// Try and see if any of the parties have your party members in it without you.
				if (partyIndex === -1) {
					console.warn('No empty temp party found for player ' + player.name);
					partyIndex = tempParties
						.filter((party) => party.length < MAX_PARTY_SIZE)
						.findIndex((party) => party.some((member) => player.partyMembers.includes(member.name)));
					console.info('Found party with members in it: ' + partyIndex);
				}

				let party: WCLSimPlayer[] = tempParties[partyIndex];

				partyMembers.forEach((partyMember) => {

					const isUndefined = party === undefined;
					const isFull = party.length === MAX_PARTY_SIZE;

					if (isUndefined || isFull) {
						return;
					}

					partyMember.partyAssigned = true;
					partyMember.partyMembers = player.partyMembers;

					party.push(partyMember);
				});
			}

			// Process the temp groups into the sim raid groups.
			tempParties.forEach((party, partyIndex) => {

				let raidParty = raid.parties[partyIndex];

				party
					.sort(sortBuffBotsLast)
					.forEach((player) => assignPlayerToParty(player, raidParty));
			});
		}

		// Process the players who didn't get assigned a group yet.
		wclPlayers
			.filter((player) => !player.partyAssigned)
			.sort(sortBuffBotsLast)
			.map((player) => {
				let raidParty = raid.parties.filter((party) => party.players.length < MAX_PARTY_SIZE)[0];
				assignPlayerToParty(player, raidParty, true);
			});

		// Insert the innervate / PI buffs into the options for the raid.
		wclPlayers
			.filter((player) => player.innervateTarget || player.powerInfusionTarget)
			.forEach((player) => {

				const target: wclSimPlayer | undefined = wclPlayers.find((wclPlayer) => wclPlayer.name === player.innervateTarget || player.name === player.powerInfusionTarget);

				if (!target) {
					console.warn('Could not find target assignment player');
					return;
				}

				const targetID: number = target.id;
				const targetRaidIndex: number | undefined = wclIDtoRaidIndex.get(targetID);

				if (!targetRaidIndex) {
					console.warn(`Could not find ${target.name} raid index!`);
					return;
				}

				if (player.isBuffBot) {
					const playerID: number = player.id;
					const playerRaidIndex: number | undefined = wclIDtoRaidIndex.get(playerID);
					const buffBot = buffBots.find((buffBot) => buffBot.raidIndex === playerRaidIndex);
					if (buffBot) {
						if (player.innervateTarget) {
							buffBot.innervateAssignment = RaidTarget.create();
							buffBot.innervateAssignment.targetIndex = targetRaidIndex
						} else if (player.powerInfusionTarget) {
							buffBot.powerInfusionAssignment = RaidTarget.create();
							buffBot.powerInfusionAssignment.targetIndex = targetRaidIndex
						}
					}
					return;
				}

				// Regular players.

				const raidParty = raid.parties.filter((party) => party.players.some((raidPlayer) => raidPlayer.name === player.name))[0];

				if (!raidParty) {
					console.warn('Could not find raiding party for player ' + player.name);
					return;
				}

				const raidPlayer = raidParty.players.find((raidPlayer) => raidPlayer.name === player.name);

				if (!raidPlayer) {
					console.warn('Could not find raid player ' + player.name + ' in raid party ' + raidParty);
					return;
				}

				if (player.innervateTarget) {
					if (raidPlayer.spec.oneofKind == 'balanceDruid') {
						raidPlayer.spec.balanceDruid.options!.innervateTarget = RaidTarget.create();
						raidPlayer.spec.balanceDruid.options!.innervateTarget.targetIndex = targetRaidIndex;
					} else if (raidPlayer.spec.oneofKind == 'feralDruid') {
						raidPlayer.spec.feralDruid.options!.innervateTarget = RaidTarget.create();
						raidPlayer.spec.feralDruid.options!.innervateTarget.targetIndex = targetRaidIndex;
					} else if (raidPlayer.spec.oneofKind == 'feralTankDruid') {
						raidPlayer.spec.feralTankDruid.options!.innervateTarget = RaidTarget.create();
						raidPlayer.spec.feralTankDruid.options!.innervateTarget.targetIndex = targetRaidIndex;
					}
				} else if (player.powerInfusionTarget) {
					// Pretty sure there is no shadow priest that has PI
				}
			});

		wclPlayers
			.filter((player) => !player.partyAssigned)
			.forEach((player) => {
				console.error(`${player.name} is not in a party!`, player);
			});

		settings.blessings = makeDefaultBlessings(numPaladins);

		this.simUI.fromProto(TypedEvent.nextEventID(), settings);
		this.simUI.setBuffBots(TypedEvent.nextEventID(), buffBots);

		if (!experimentalGenerateParties) {
			const rateLimitEnd: wclRateLimitData = await this.getRateLimit(token);
			console.debug(`Rate Limit resets in ${rateLimitEnd.pointsResetIn} seconds.`);
		}

		this.close();
	}
}

class WCLSimPlayer implements wclSimPlayer {
	public gear: wclGear[];
	public icon: string;
	public id: number;
	public name: string;
	public type: string;
	public talents: wclTalents[];
	public wclSpec: string;

	public partyAssigned: boolean = false;
	public partyFound: boolean = false;
	public partyMembers: string[] = [];

	public isPartyBuffer: boolean = false;
	public isBuffBot: boolean = false;
	public sortPriority: number = 99;

	public innervateTarget: string | undefined;
	public powerInfusionTarget: string | undefined;

	private simUI: RaidSimUI;
	private spec: Spec;
	private specType: string;
	private faction: Faction;

	constructor(data: wclPlayer, simUI: RaidSimUI, faction: Faction = Faction.Unknown) {
		this.simUI = simUI;

		this.name = data.name;
		this.gear = data.gear;
		this.icon = data.icon;
		this.id = data.id;
		this.type = data.type;
		this.talents = data.talents;
		this.wclSpec = data.icon.split('-')[1];
		this.faction = faction;

		// Prot Paladin's occasionally have a specType of 'Protection' instead of 'Justicar'?
		if (this.type === 'Paladin' && this.wclSpec === 'Protection') {
			this.wclSpec = 'Justicar';
		}

		this.spec = specNames[this.wclSpec];
		this.specType = this.wclSpec + this.type;

		this.isBuffBot = this.spec === undefined;

		this.isPartyBuffer = this.getPartyAuraIds().length > 1;

		this.sortPriority = specSortPriority[this.wclSpec] ?? 99;
	}

	public getPlayer(): PlayerProto | undefined {

		if (this.isBuffBot) {
			return;
		}

		let player = PlayerProto.create();

		const specFuncs = specTypeFunctions[this.spec];

		const matchingPreset = this.getMatchingPreset();

		if (matchingPreset === undefined) {
			console.error('Could not find matching preset for non buff bot', {
				'name': this.name,
				'spec': this.spec,
				'type': this.type,
				'talents': this.talents,
			});
			return;
		}

		player = withSpecProto(this.spec, player, matchingPreset.rotation, specFuncs.talentsCreate(), matchingPreset.specOptions);

		player.talentsString = matchingPreset.talents;
		player.consumes = matchingPreset.consumes;

		player.name = this.name;
		player.class = nameToClass(this.type);
		player.equipment = this.getEquipment();
		player.race = matchingPreset.defaultFactionRaces[this.faction];

		Player.applySharedDefaultsToProto(player);

		return player;
	}

	public getBuffBot(): BuffBot | undefined {

		if (!this.isBuffBot) {
			return;
		}

		const botID = buffBotNames[this.specType];

		if (botID == null) {
			console.error('Buff Bot Spec not implemented: ', this.specType);
			return;
		}

		const bot = BuffBot.create();
		bot.id = botID;
		bot.raidIndex = -1; // Default it for now. // numPlayers

		return bot;
	}

	public getPartyAuraIds() {

		const allSpecClassAuras: any = {
			'Paladin': [
				19746, // Concentration Aura
				27149, // Devotion Aura,
				27150, // Retribution Aura
			],
			'Warrior': [
				2048,  // Battle Shout
				469, // Commanding Shout
			],
			'Warlock': [
				27268, // Pet Imp: Blood Pact
				18696, // Improved Imp: Blood Pact
			],
		};

		// Reused for the plethora of Feral specs.
		const feralDruidSpecAuras = [
			24932, // Improved Leader of the Pack // at least 0,32,0
			// 17007, // Leader of the Pack // at least 0,31,0
		];

		// TODO: Could additionally filter out buff IDs based on minimum req talent strings?
		const specSpecificAuras: any = {
			'RetributionPaladin': [
				20092, // Improved Retribution Aura // at least 0,0,16
				20218, // Sanctity Aura // at least 0,0,21
				31870, // Improved Sanctity Aura // at least 0,0,22
			],
			'GuardianDruid': [...feralDruidSpecAuras],
			'WardenDruid': [...feralDruidSpecAuras],
			'FeralDruid': [...feralDruidSpecAuras],
			'BalanceDruid': [
				24907, // Moonkin Aura // at least 31,0,0
			],
			'RestorationDruid': [
				34123, // Tree of Life // at least 0,0,41
			],
			'MarksmanHunter': [
				27066, // Trueshot Aura // at least 0,32,0
			],
			'EnhancementShaman': [
				30811, // Unleashed Rage // at least 0,36,0
			],
			// 'ElementalShaman': [] // Totem buffs do not show up in logs. Leaving for future reference.
		};

		const consumableAuras = [
			351355, // Greater Drums of Battle
		];

		const classAuras = allSpecClassAuras[this.type] ?? [];
		const specAuras = specSpecificAuras[this.specType] ?? [];

		const reliableAuras = [
			...specAuras, ...classAuras, ...consumableAuras,
		];

		if (this.type === 'Shaman') {
			// Shamans get moved around a lot, so Heroism isn't a good reference for what group they are in.
			return [
				...reliableAuras,
				32182, // Heroism
			];
		}
		return reliableAuras;
	}

	private getMatchingPreset(): PresetSpecSettings<Spec> {
		const matchingPresets = playerPresets.filter((preset) => preset.spec === this.spec);
		let presetIdx = 0;

		if (matchingPresets && matchingPresets.length > 1) {
			let distance = 100;
			// Search talents and find the preset that the players talents most closely match.
			matchingPresets.forEach((preset, i) => {
				let presetTalents = [0, 0, 0];
				let talentIdx = 0;
				// First sum up the number of talents per tree for preset.
				Array.from(preset.talents).forEach((v) => {
					if (v == '-') {
						talentIdx++;
						return;
					}
					presetTalents[talentIdx] += parseInt(v);
				});

				// Diff the distance to the preset.

				const newDistance = presetTalents.reduce((acc, v, i) => acc += Math.abs(this.talents[i]?.guid - presetTalents[i]), 0);

				// If this is the best distance, assign this preset.
				if (newDistance < distance) {
					presetIdx = i;
					distance = newDistance;
				}
			});
		}
		return matchingPresets[presetIdx];
	}

	private getEquipment(): EquipmentSpec {
		let equipment = EquipmentSpec.create();
		equipment.items = new Array<ItemSpec>();

		this.gear.forEach((gear) => {
			const item = ItemSpec.create();
			item.id = gear.id;
			const dbEnchant = this.simUI.sim.getEnchantFlexible(gear.permanentEnchant);
			item.enchant = dbEnchant
				? dbEnchant.id
				: 0;
			if (gear.gems) {
				item.gems = new Array<number>();
				gear.gems.forEach((gemInfo) => item.gems.push(gemInfo.id));
			}
			equipment!.items.push(item);
		});
		return equipment;
	}

}


// Maps WCL spec to sorting priority for party makeup checks. Lower the number, the more likely the query will be successful.
const specSortPriority: Record<string, number> = {
	'Warden': 0,
	'Guardian': 1,
	'Feral': 2,
	'Balance': 3,
	'Justicar': 4,
	'Retribution': 5,
	'Fury': 6,
	'Arms': 7,
	'Protection': 8,
	'Enhancement': 9,
	'Destruction': 10,
	'Affliction': 11,
	'Demonology': 12,
	'Marksman': 13,
};

// Maps WCL spec names to internal Spec enum.
const specNames: Record<string, Spec> = {
	'Balance': Spec.SpecBalanceDruid,
	'Elemental': Spec.SpecElementalShaman,
	'Enhancement': Spec.SpecEnhancementShaman,
	'Feral': Spec.SpecFeralDruid,
	'Warden': Spec.SpecFeralTankDruid,
	'Guardian': Spec.SpecFeralTankDruid,
	'Survival': Spec.SpecHunter,
	'BeastMastery': Spec.SpecHunter,
	'Arcane': Spec.SpecMage,
	'Fire': Spec.SpecMage,
	'Frost': Spec.SpecMage,
	'Assassination': Spec.SpecRogue,
	'Combat': Spec.SpecRogue,
	'Retribution': Spec.SpecRetributionPaladin,
	'Justicar': Spec.SpecProtectionPaladin,
	'Shadow': Spec.SpecShadowPriest,
	'Smite': Spec.SpecSmitePriest,
	'Destruction': Spec.SpecWarlock,
	'Affliction': Spec.SpecWarlock,
	'Demonology': Spec.SpecWarlock,
	'Arms': Spec.SpecWarrior,
	'Fury': Spec.SpecWarrior,
	'Champion': Spec.SpecWarrior,
	'Warrior': Spec.SpecWarrior,
	'Gladiator': Spec.SpecWarrior,
	'Protection': Spec.SpecProtectionWarrior,
};

// Maps WCL spec+type to internal buff bot IDs.
const buffBotNames: Record<string, string> = {
	// Healers
	'HolyPaladin': 'Paladin',
	'HolyPriest': 'Holy Priest',
	'DisciplinePriest': 'Divine Spirit Priest',
	'RestorationDruid': 'Resto Druid',
	'DreamstateDruid': 'Resto Druid',
	'RestorationShaman': 'Resto Shaman',
};

interface wclBuffCastsData {
	name: string;
	targets: {
		name: string;
		type: string;
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

// Typed interface for WoWSimPlayer class
interface wclSimPlayer extends wclPlayer {
	wclSpec: string;
	partyAssigned: boolean;
	isPartyBuffer: boolean;
	partyMembers: string[];
	isBuffBot: boolean;
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
