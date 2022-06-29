import { Class } from '/tbc/core/proto/common.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { BuffBot as BuffBotProto } from '/tbc/core/proto/ui.js';
import { classColors } from '/tbc/core/proto_utils/utils.js';
import { emptyRaidTarget } from '/tbc/core/proto_utils/utils.js';
import { specToClass } from '/tbc/core/proto_utils/utils.js';
import { Sim } from '/tbc/core/sim.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { BuffBotSettings, buffBotPresets } from './presets.js';

export const NO_ASSIGNMENT = -1;

// Represents a buff bot in a raid.
export class BuffBot {
	settings: BuffBotSettings;
	spec: Spec = 0;
	name: string = '';

	private raidIndex: number = NO_ASSIGNMENT;
	private innervateAssignment: RaidTarget = emptyRaidTarget();
	private powerInfusionAssignment: RaidTarget = emptyRaidTarget();

	readonly raidIndexChangeEmitter = new TypedEvent<void>();
	readonly innervateAssignmentChangeEmitter = new TypedEvent<void>();
	readonly powerInfusionAssignmentChangeEmitter = new TypedEvent<void>();
	readonly changeEmitter = new TypedEvent<void>();

	private readonly sim: Sim;

	constructor(id: string, sim: Sim) {
		const settings = buffBotPresets.find(preset => preset.buffBotId == id);
		if (!settings) {
			throw new Error('No buff bot config with id \'' + id + '\'!');
		}

		this.sim = sim;
		this.settings = settings;
		this.updateSettings();

		[
			this.raidIndexChangeEmitter,
			this.innervateAssignmentChangeEmitter,
			this.powerInfusionAssignmentChangeEmitter,
		].forEach(emitter => emitter.on(eventID => this.changeEmitter.emit(eventID)));

		this.changeEmitter.on(eventID => sim.raid.getParty(this.getPartyIndex()).changeEmitter.emit(eventID));
	}

	private updateSettings() {
		this.spec = this.settings.spec;
		this.name = this.settings.name;
	}

	getLabel(): string {
		return `${this.name} (#${this.getRaidIndex() + 1})`;
	}

	getClass(): Class {
		return specToClass[this.settings.spec];
	}

	getClassColor(): string {
		return classColors[this.getClass()];
	}

	getRaidIndex(): number {
		return this.raidIndex;
	}
	setRaidIndex(eventID: EventID, newRaidIndex: number) {
		if (newRaidIndex != this.raidIndex) {
			this.raidIndex = newRaidIndex;
			TypedEvent.freezeAllAndDo(() => {
				this.raidIndexChangeEmitter.emit(eventID);
				this.sim.raid.compChangeEmitter.emit(eventID);
			});
		}
	}

	getPartyIndex(): number {
		return Math.floor(this.getRaidIndex() / 5);
	}

	getInnervateAssignment(): RaidTarget {
		// Defensive copy.
		return RaidTarget.clone(this.innervateAssignment);
	}
	setInnervateAssignment(eventID: EventID, newInnervateAssignment: RaidTarget) {
		if (RaidTarget.equals(newInnervateAssignment, this.innervateAssignment))
			return;

		// Defensive copy.
		this.innervateAssignment = RaidTarget.clone(newInnervateAssignment);
		this.innervateAssignmentChangeEmitter.emit(eventID);
	}

	getPowerInfusionAssignment(): RaidTarget {
		// Defensive copy.
		return RaidTarget.clone(this.powerInfusionAssignment);
	}
	setPowerInfusionAssignment(eventID: EventID, newPowerInfusionAssignment: RaidTarget) {
		if (RaidTarget.equals(newPowerInfusionAssignment, this.powerInfusionAssignment))
			return;

		// Defensive copy.
		this.powerInfusionAssignment = RaidTarget.clone(newPowerInfusionAssignment);
		this.powerInfusionAssignmentChangeEmitter.emit(eventID);
	}

	toProto(): BuffBotProto {
		return BuffBotProto.create({
			id: this.settings.buffBotId,
			raidIndex: this.getRaidIndex(),
			innervateAssignment: this.getInnervateAssignment(),
			powerInfusionAssignment: this.getPowerInfusionAssignment(),
		});
	}

	fromProto(eventID: EventID, proto: BuffBotProto) {
		const settings = buffBotPresets.find(preset => preset.buffBotId == proto.id);
		if (!settings) {
			throw new Error('No buff bot config with id \'' + proto.id + '\'!');
		}
		this.settings = settings;
		this.updateSettings();
		TypedEvent.freezeAllAndDo(() => {
			this.setRaidIndex(eventID, proto.raidIndex);
			this.setInnervateAssignment(eventID, proto.innervateAssignment || emptyRaidTarget());
			this.setPowerInfusionAssignment(eventID, proto.powerInfusionAssignment || emptyRaidTarget());
		});
	}

	clone(eventID: EventID): BuffBot {
		const newBuffBot = new BuffBot(this.settings.buffBotId, this.sim);
		newBuffBot.fromProto(eventID, this.toProto());
		return newBuffBot;
	}
}
