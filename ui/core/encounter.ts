import { Encounter as EncounterProto } from '/wotlk/core/proto/common.js';
import { MobType } from '/wotlk/core/proto/common.js';
import { Stat } from '/wotlk/core/proto/common.js';
import { Target as TargetProto } from '/wotlk/core/proto/common.js';
import { PresetEncounter } from '/wotlk/core/proto/api.js';
import { PresetTarget } from '/wotlk/core/proto/api.js';
import { Target } from '/wotlk/core/target.js';
import { Stats } from '/wotlk/core/proto_utils/stats.js';

import { Sim } from './sim.js';
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
	private targets: Array<Target>;

	readonly targetsChangeEmitter = new TypedEvent<void>();
	readonly durationChangeEmitter = new TypedEvent<void>();
	readonly executeProportionChangeEmitter = new TypedEvent<void>();

	// Emits when any of the above emitters emit.
	readonly changeEmitter = new TypedEvent<void>();

	constructor(sim: Sim) {
		this.sim = sim;
		this.targets = [Target.fromDefaults(TypedEvent.nextEventID(), sim)];

		[
			this.targetsChangeEmitter,
			this.durationChangeEmitter,
			this.executeProportionChangeEmitter,
		].forEach(emitter => emitter.on(eventID => this.changeEmitter.emit(eventID)));
	}

	get primaryTarget(): Target {
		return this.targets[0];
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

	getNumTargets(): number {
		return this.targets.length;
	}

	getTargets(): Array<Target> {
		return this.targets.slice();
	}
	setTargets(eventID: EventID, newTargets: Array<Target>) {
		TypedEvent.freezeAllAndDo(() => {
			if (newTargets.length == 0) {
				newTargets = [Target.fromDefaults(eventID, this.sim)];
			}
			if (newTargets.length == this.targets.length && newTargets.every((target, i) => TargetProto.equals(target.toProto(), this.targets[i].toProto()))) {
				return;
			}

			this.targets = newTargets;
			this.targetsChangeEmitter.emit(eventID);
		});
	}

	matchesPreset(preset: PresetEncounter): boolean {
		return preset.targets.length == this.targets.length && this.targets.every((t, i) => t.matchesPreset(preset.targets[i]));
	}

	applyPreset(eventID: EventID, preset: PresetEncounter) {
		TypedEvent.freezeAllAndDo(() => {
			let newTargets = this.targets.slice(0, preset.targets.length);
			while (newTargets.length < preset.targets.length) {
				newTargets.push(new Target(this.sim));
			}

			newTargets.forEach((nt, i) => nt.applyPreset(eventID, preset.targets[i]));
			this.setTargets(eventID, newTargets);
		});
	}

	toProto(): EncounterProto {
		return EncounterProto.create({
			duration: this.duration,
			durationVariation: this.durationVariation,
			executeProportion20: this.executeProportion20,
			executeProportion25: this.executeProportion25,
			executeProportion35: this.executeProportion35,
			useHealth: this.useHealth,
			targets: this.targets.map(target => target.toProto()),
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

			if (proto.targets.length > 0) {
				this.setTargets(eventID, proto.targets.map(targetProto => {
					const target = new Target(this.sim);
					target.fromProto(eventID, targetProto);
					return target;
				}));
			} else {
				this.setTargets(eventID, [Target.fromDefaults(eventID, this.sim)]);
			}
		});
	}

	applyDefaults(eventID: EventID) {
		this.fromProto(eventID, EncounterProto.create({
			duration: 180,
			durationVariation: 5,
			executeProportion20: 0.2,
			executeProportion25: 0.25,
			executeProportion35: 0.35,
			targets: [Target.defaultProto()],
		}));
	}
}
