import {
	Class,
	Debuffs,
	RaidBuffs,
	UnitReference,
	UnitReference_Type as UnitType,
} from './proto/common.js';
import { Raid as RaidProto } from './proto/api.js';

import { Party, MAX_PARTY_SIZE } from './party.js';
import { Player } from './player.js';
import { EventID, TypedEvent } from './typed_event.js';
import { Sim } from './sim.js';
import { sum } from './utils.js';

export const MAX_NUM_PARTIES = 8;

// Manages all the settings for a single Raid.
export class Raid {
	private buffs: RaidBuffs = RaidBuffs.create();
	private debuffs: Debuffs = Debuffs.create();
	private tanks: Array<UnitReference> = [];
	private targetDummies: number = 0;
	private numActiveParties: number = 5;

	// Emits when a raid member is added/removed/moved.
	readonly compChangeEmitter = new TypedEvent<void>();

	readonly buffsChangeEmitter = new TypedEvent<void>();
	readonly debuffsChangeEmitter = new TypedEvent<void>();
	readonly tanksChangeEmitter = new TypedEvent<void>();
	readonly targetDummiesChangeEmitter = new TypedEvent<void>();
	readonly numActivePartiesChangeEmitter = new TypedEvent<void>();

	// Emits when anything in the raid changes.
	readonly changeEmitter: TypedEvent<void>;

	// Should always hold exactly MAX_NUM_PARTIES elements.
	private parties: Array<Party>;

	// Cached return value for getActivePlayers().
	private activePlayers: Array<Player<any>>;

	readonly sim: Sim;

	constructor(sim: Sim) {
		this.sim = sim;

		this.parties = [...Array(MAX_NUM_PARTIES).keys()].map(i => {
			const newParty = new Party(this, sim);
			newParty.compChangeEmitter.on(eventID => this.compChangeEmitter.emit(eventID));
			newParty.changeEmitter.on(eventID => this.changeEmitter.emit(eventID));
			return newParty;
		});
		this.activePlayers = [];

		this.numActivePartiesChangeEmitter.on(eventID => this.compChangeEmitter.emit(eventID));

		this.changeEmitter = TypedEvent.onAny([
			this.compChangeEmitter,
			this.buffsChangeEmitter,
			this.debuffsChangeEmitter,
			this.tanksChangeEmitter,
			this.targetDummiesChangeEmitter,
		], 'RaidChange');

		this.changeEmitter.on(() => {
			this.activePlayers = [];
		});
	}

	size(): number {
		return sum(this.parties.map(party => party.size()));
	}

	isEmpty(): boolean {
		return this.size() == 0;
	}

	getParties(): Array<Party> {
		// Make defensive copy.
		return this.parties.slice();
	}

	getParty(index: number): Party {
		return this.parties[index];
	}

	getPlayers(): Array<Player<any> | null> {
		return this.parties.map(party => party.getPlayers()).flat();
	}

	getPlayer(index: number): Player<any> | null {
		if (index === -1) return null;

		const party = this.parties[Math.floor(index / MAX_PARTY_SIZE)];
		return party.getPlayer(index % MAX_PARTY_SIZE);
	}

	getPlayerFromUnitReference(raidTarget: UnitReference|undefined, contextPlayer?: Player<any>|null): Player<any> | null {
		if (!raidTarget || raidTarget.type == UnitType.Unknown) {
			return null;
		} else if (raidTarget.type == UnitType.Player) {
			return this.getPlayer(raidTarget.index);
		} else if (raidTarget.type == UnitType.Self) {
			return contextPlayer || null;
		} else {
			return null;
		}
	}

	setPlayer(eventID: EventID, index: number, newPlayer: Player<any> | null) {
		const party = this.parties[Math.floor(index / MAX_PARTY_SIZE)];
		party.setPlayer(eventID, index % MAX_PARTY_SIZE, newPlayer);
	}

	getClassCount(playerClass: Class) {
		return this.getPlayers().filter(player => player != null && player.getClass() == playerClass).length;
	}

	getBuffs(): RaidBuffs {
		// Make a defensive copy
		return RaidBuffs.clone(this.buffs);
	}

	setBuffs(eventID: EventID, newBuffs: RaidBuffs) {
		if (RaidBuffs.equals(this.buffs, newBuffs))
			return;


		// Make a defensive copy
		this.buffs = RaidBuffs.clone(newBuffs);

		// Special handle ToW since it crosses buffs/debuffs.
		if (this.debuffs.totemOfWrath != this.buffs.totemOfWrath) {
			var newDebuff = Debuffs.clone(this.debuffs);
			newDebuff.totemOfWrath = this.buffs.totemOfWrath;
			this.setDebuffs(eventID, newDebuff);
		}
		this.buffsChangeEmitter.emit(eventID);
	}

	getDebuffs(): Debuffs {
		// Make a defensive copy
		return Debuffs.clone(this.debuffs);
	}

	setDebuffs(eventID: EventID, newDebuffs: Debuffs) {
		if (Debuffs.equals(this.debuffs, newDebuffs))
			return;

		// Make a defensive copy
		this.debuffs = Debuffs.clone(newDebuffs);

		// Special handle ToW since it crosses buffs/debuffs.
		if (this.debuffs.totemOfWrath != this.buffs.totemOfWrath) {
			var newBuffs = RaidBuffs.clone(this.buffs);
			newBuffs.totemOfWrath = this.debuffs.totemOfWrath;
			this.setBuffs(eventID, newBuffs);
		}
		this.debuffsChangeEmitter.emit(eventID);
	}

	getTanks(): Array<UnitReference> {
		// Make a defensive copy
		return this.tanks.map(tank => UnitReference.clone(tank));
	}

	setTanks(eventID: EventID, newTanks: Array<UnitReference>) {
		if (this.tanks.length == newTanks.length && this.tanks.every((tank, i) => UnitReference.equals(tank, newTanks[i])))
			return;

		// Make a defensive copy
		this.tanks = newTanks.map(tank => UnitReference.clone(tank));
		this.tanksChangeEmitter.emit(eventID);
	}

	getTargetDummies(): number {
		return this.targetDummies;
	}

	setTargetDummies(eventID: EventID, newTargetDummies: number) {
		if (this.targetDummies == newTargetDummies)
			return;

		this.targetDummies = newTargetDummies;
		this.targetDummiesChangeEmitter.emit(eventID);
	}

	getNumActiveParties(): number {
		return this.numActiveParties;
	}
	setNumActiveParties(eventID: EventID, newNumActiveParties: number) {
		if (newNumActiveParties != this.numActiveParties && newNumActiveParties > 0) {
			this.numActiveParties = newNumActiveParties;
			this.numActivePartiesChangeEmitter.emit(eventID);
		}
	}
	getActivePlayers(): Array<Player<any>> {
		if (this.activePlayers.length == 0) {
			const activeParties = this.getParties().filter((party, i) => i < this.numActiveParties);
			this.activePlayers = activeParties
				.map(party => party.getPlayers())
				.flat()
				.filter(player => player != null) as Array<Player<any>>;
		}
		return this.activePlayers;
	}

	toProto(forExport?: boolean, forSimming?: boolean): RaidProto {
		return RaidProto.create({
			parties: this.parties.map(party => party.toProto(forExport, forSimming)),
			buffs: this.getBuffs(),
			debuffs: this.getDebuffs(),
			tanks: this.getTanks(),
			targetDummies: this.getTargetDummies(),
			numActiveParties: this.getNumActiveParties(),
		});
	}

	fromProto(eventID: EventID, proto: RaidProto) {
		TypedEvent.freezeAllAndDo(() => {
			this.setBuffs(eventID, proto.buffs || RaidBuffs.create());
			this.setDebuffs(eventID, proto.debuffs || Debuffs.create());
			this.setTanks(eventID, proto.tanks);
			this.setTargetDummies(eventID, proto.targetDummies);
			this.setNumActiveParties(eventID, proto.numActiveParties || 5);

			for (let i = 0; i < MAX_NUM_PARTIES; i++) {
				if (proto.parties[i]) {
					this.parties[i].fromProto(eventID, proto.parties[i]);
				} else {
					this.parties[i].clear(eventID);
				}
			}
		});
	}

	clear(eventID: EventID) {
		TypedEvent.freezeAllAndDo(() => {
			for (let i = 0; i < MAX_NUM_PARTIES; i++) {
				this.parties[i].clear(eventID);
			}
		});
	}
}
