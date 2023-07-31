import {
	Encounter as EncounterProto,
	MobType,
	SpellSchool,
	Stat,
	Target as TargetProto,
	TargetInput,
	PresetEncounter,
	PresetTarget,
} from './proto/common.js';
import { Stats } from './proto_utils/stats.js';
import * as Mechanics from './constants/mechanics.js';

import { Sim } from './sim.js';
import { UnitMetadataList } from './player.js';
import { EventID, TypedEvent } from './typed_event.js';

// Manages all the settings for an Encounter.
export class Encounter {
	readonly sim: Sim;

	private duration: number = 180;
	private durationVariation: number = 5;
	private executeProportion20: number = 0.2;
	private executeProportion25: number = 0.25;
	private executeProportion35: number = 0.35;
	private useHealth: boolean = false;
	targets: Array<TargetProto>;
	targetsMetadata: UnitMetadataList;

	readonly targetsChangeEmitter = new TypedEvent<void>();
	readonly durationChangeEmitter = new TypedEvent<void>();
	readonly executeProportionChangeEmitter = new TypedEvent<void>();

	// Emits when any of the above emitters emit.
	readonly changeEmitter = new TypedEvent<void>();

	constructor(sim: Sim) {
		this.sim = sim;
		this.targets = [Encounter.defaultTargetProto()];
		this.targetsMetadata = new UnitMetadataList();

		[
			this.targetsChangeEmitter,
			this.durationChangeEmitter,
			this.executeProportionChangeEmitter,
		].forEach(emitter => emitter.on(eventID => this.changeEmitter.emit(eventID)));
	}

	get primaryTarget(): TargetProto {
		return TargetProto.clone(this.targets[0]);
	}

	getDurationVariation(): number {
		return this.durationVariation;
	}
	setDurationVariation(eventID: EventID, newDuration: number) {
		if (newDuration == this.durationVariation)
			return;

		this.durationVariation = newDuration;
		this.durationChangeEmitter.emit(eventID);
	}

	getDuration(): number {
		return this.duration;
	}
	setDuration(eventID: EventID, newDuration: number) {
		if (newDuration == this.duration)
			return;

		this.duration = newDuration;
		this.durationChangeEmitter.emit(eventID);
	}

	getExecuteProportion20(): number {
		return this.executeProportion20;
	}
	setExecuteProportion20(eventID: EventID, newExecuteProportion20: number) {
		if (newExecuteProportion20 == this.executeProportion20)
			return;

		this.executeProportion20 = newExecuteProportion20;
		this.executeProportionChangeEmitter.emit(eventID);
	}
	getExecuteProportion25(): number {
		return this.executeProportion25;
	}
	setExecuteProportion25(eventID: EventID, newExecuteProportion25: number) {
		if (newExecuteProportion25 == this.executeProportion25)
			return;

		this.executeProportion25 = newExecuteProportion25;
		this.executeProportionChangeEmitter.emit(eventID);
	}
	getExecuteProportion35(): number {
		return this.executeProportion35;
	}
	setExecuteProportion35(eventID: EventID, newExecuteProportion35: number) {
		if (newExecuteProportion35 == this.executeProportion35)
			return;

		this.executeProportion35 = newExecuteProportion35;
		this.executeProportionChangeEmitter.emit(eventID);
	}

	getUseHealth(): boolean {
		return this.useHealth;
	}
	setUseHealth(eventID: EventID, newUseHealth: boolean) {
		if (newUseHealth == this.useHealth)
			return;

		this.useHealth = newUseHealth;
		this.durationChangeEmitter.emit(eventID);
		this.executeProportionChangeEmitter.emit(eventID);
	}

	matchesPreset(preset: PresetEncounter): boolean {
		return preset.targets.length == this.targets.length && this.targets.every((t, i) => TargetProto.equals(t, preset.targets[i].target));
	}

	applyPreset(eventID: EventID, preset: PresetEncounter) {
		this.targets = preset.targets.map(presetTarget => presetTarget.target || TargetProto.create());
		this.targetsChangeEmitter.emit(eventID);
	}

	applyPresetTarget(eventID: EventID, preset: PresetTarget, index: number) {
		this.targets[index] = preset.target || TargetProto.create();
		this.targetsChangeEmitter.emit(eventID);
	}

	toProto(): EncounterProto {
		return EncounterProto.create({
			duration: this.duration,
			durationVariation: this.durationVariation,
			executeProportion20: this.executeProportion20,
			executeProportion25: this.executeProportion25,
			executeProportion35: this.executeProportion35,
			useHealth: this.useHealth,
			targets: this.targets,
		});
	}

	fromProto(eventID: EventID, proto: EncounterProto) {
		TypedEvent.freezeAllAndDo(() => {
			this.setDuration(eventID, proto.duration);
			this.setDurationVariation(eventID, proto.durationVariation);
			this.setExecuteProportion20(eventID, proto.executeProportion20);
			this.setExecuteProportion25(eventID, proto.executeProportion25);
			this.setExecuteProportion35(eventID, proto.executeProportion35);
			this.setUseHealth(eventID, proto.useHealth);
			this.targets = proto.targets;
			this.targetsChangeEmitter.emit(eventID);
		});
	}

	applyDefaults(eventID: EventID) {
		this.fromProto(eventID, EncounterProto.create({
			duration: 180,
			durationVariation: 5,
			executeProportion20: 0.2,
			executeProportion25: 0.25,
			executeProportion35: 0.35,
			targets: [Encounter.defaultTargetProto()],
		}));
	}

	static defaultTargetProto(): TargetProto {
		return TargetProto.create({
			level: Mechanics.BOSS_LEVEL,
			mobType: MobType.MobTypeGiant,
			tankIndex: 0,
			swingSpeed: 1.5,
			minBaseDamage: 65000,
			dualWield: false,
			dualWieldPenalty: false,
			suppressDodge: false,
			parryHaste: true,
			spellSchool: SpellSchool.SpellSchoolPhysical,
			stats: Stats.fromMap({
				[Stat.StatArmor]: 10643,
				[Stat.StatAttackPower]: 805,
				[Stat.StatBlockValue]: 76,
			}).asArray(),
			targetInputs: new Array<TargetInput>(0),
		});
	}
}
