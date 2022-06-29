import { MobType } from '/tbc/core/proto/common.js';
import { SpellSchool } from '/tbc/core/proto/common.js';
import { Stat } from '/tbc/core/proto/common.js';
import { Target as TargetProto } from '/tbc/core/proto/common.js';
import { PresetTarget } from '/tbc/core/proto/api.js';
import { Stats } from '/tbc/core/proto_utils/stats.js';

import * as Mechanics from '/tbc/core/constants/mechanics.js';

import { Listener } from './typed_event.js';
import { Sim } from './sim.js';
import { EventID, TypedEvent } from './typed_event.js';
import { sum } from './utils.js';
import { wait } from './utils.js';

// Manages all the settings for a single Target.
export class Target {
	readonly sim: Sim;

	private id: number = 0;
	private name: string = '';
	private level: number = Mechanics.BOSS_LEVEL;
	private mobType: MobType = MobType.MobTypeDemon;
	private tankIndex: number = 0;
	private stats: Stats = new Stats();

	private swingSpeed: number = 0;
	private minBaseDamage: number = 0;
	private dualWield: boolean = false;
	private dualWieldPenalty: boolean = false;
	private canCrush: boolean = true;
	private suppressDodge: boolean = false;
	private parryHaste: boolean = true;
	private spellSchool: SpellSchool = SpellSchool.SpellSchoolPhysical;

	readonly idChangeEmitter = new TypedEvent<void>();
	readonly nameChangeEmitter = new TypedEvent<void>();
	readonly levelChangeEmitter = new TypedEvent<void>();
	readonly mobTypeChangeEmitter = new TypedEvent<void>();
	readonly propChangeEmitter = new TypedEvent<void>();
	readonly statsChangeEmitter = new TypedEvent<void>();

	// Emits when any of the above emitters emit.
	readonly changeEmitter = new TypedEvent<void>();

	constructor(sim: Sim) {
		this.sim = sim;

		[
			this.idChangeEmitter,
			this.nameChangeEmitter,
			this.levelChangeEmitter,
			this.mobTypeChangeEmitter,
			this.propChangeEmitter,
			this.statsChangeEmitter,
		].forEach(emitter => emitter.on(eventID => this.changeEmitter.emit(eventID)));

		this.changeEmitter.on(eventID => this.sim.encounter?.changeEmitter.emit(eventID));
	}

	getId(): number {
		return this.id;
	}

	setId(eventID: EventID, newId: number) {
		if (newId == this.id)
			return;

		this.id = newId;
		this.idChangeEmitter.emit(eventID);
	}

	getName(): string {
		return this.name;
	}

	setName(eventID: EventID, newName: string) {
		if (newName == this.name)
			return;

		this.name = newName;
		this.nameChangeEmitter.emit(eventID);
	}

	getLevel(): number {
		return this.level;
	}

	setLevel(eventID: EventID, newLevel: number) {
		if (newLevel == this.level)
			return;

		this.level = newLevel;
		this.levelChangeEmitter.emit(eventID);
	}

	getMobType(): MobType {
		return this.mobType;
	}

	setMobType(eventID: EventID, newMobType: MobType) {
		if (newMobType == this.mobType)
			return;

		this.mobType = newMobType;
		this.mobTypeChangeEmitter.emit(eventID);
	}

	getTankIndex(): number {
		return this.tankIndex;
	}

	setTankIndex(eventID: EventID, newTankIndex: number) {
		if (newTankIndex == this.tankIndex)
			return;

		this.tankIndex = newTankIndex;
		this.propChangeEmitter.emit(eventID);
	}

	getSwingSpeed(): number {
		return this.swingSpeed;
	}

	setSwingSpeed(eventID: EventID, newSwingSpeed: number) {
		if (newSwingSpeed == this.swingSpeed)
			return;

		this.swingSpeed = newSwingSpeed;
		this.propChangeEmitter.emit(eventID);
	}

	getMinBaseDamage(): number {
		return this.minBaseDamage;
	}

	setMinBaseDamage(eventID: EventID, newMinBaseDamage: number) {
		if (newMinBaseDamage == this.minBaseDamage)
			return;

		this.minBaseDamage = newMinBaseDamage;
		this.propChangeEmitter.emit(eventID);
	}

	getDualWield(): boolean {
		return this.dualWield;
	}

	setDualWield(eventID: EventID, newDualWield: boolean) {
		if (newDualWield == this.dualWield)
			return;

		this.dualWield = newDualWield;
		this.propChangeEmitter.emit(eventID);
	}

	getDualWieldPenalty(): boolean {
		return this.dualWieldPenalty;
	}

	setDualWieldPenalty(eventID: EventID, newDualWieldPenalty: boolean) {
		if (newDualWieldPenalty == this.dualWieldPenalty)
			return;

		this.dualWieldPenalty = newDualWieldPenalty;
		this.propChangeEmitter.emit(eventID);
	}

	getCanCrush(): boolean {
		return this.canCrush;
	}

	setCanCrush(eventID: EventID, newCanCrush: boolean) {
		if (newCanCrush == this.canCrush)
			return;

		this.canCrush = newCanCrush;
		this.propChangeEmitter.emit(eventID);
	}

	getSuppressDodge(): boolean {
		return this.suppressDodge;
	}

	setSuppressDodge(eventID: EventID, newSuppressDodge: boolean) {
		if (newSuppressDodge == this.suppressDodge)
			return;

		this.suppressDodge = newSuppressDodge;
		this.propChangeEmitter.emit(eventID);
	}

	getParryHaste(): boolean {
		return this.parryHaste;
	}

	setParryHaste(eventID: EventID, newParryHaste: boolean) {
		if (newParryHaste == this.parryHaste)
			return;

		this.parryHaste = newParryHaste;
		this.propChangeEmitter.emit(eventID);
	}

	getSpellSchool(): SpellSchool {
		return this.spellSchool;
	}

	setSpellSchool(eventID: EventID, newSpellSchool: SpellSchool) {
		if (newSpellSchool == this.spellSchool)
			return;

		this.spellSchool = newSpellSchool;
		this.propChangeEmitter.emit(eventID);
	}

	getStats(): Stats {
		return this.stats;
	}

	setStats(eventID: EventID, newStats: Stats) {
		if (newStats.equals(this.stats))
			return;

		this.stats = newStats;
		this.statsChangeEmitter.emit(eventID);
	}

	matchesPreset(preset: PresetTarget): boolean {
		return TargetProto.equals(this.toProto(), preset.target);
	}

	applyPreset(eventID: EventID, preset: PresetTarget) {
		this.fromProto(eventID, preset.target || TargetProto.create());
	}

	toProto(): TargetProto {
		return TargetProto.create({
			id: this.getId(),
			name: this.getName(),
			level: this.getLevel(),
			mobType: this.getMobType(),
			tankIndex: this.getTankIndex(),
			swingSpeed: this.getSwingSpeed(),
			minBaseDamage: this.getMinBaseDamage(),
			dualWield: this.getDualWield(),
			dualWieldPenalty: this.getDualWieldPenalty(),
			canCrush: this.getCanCrush(),
			suppressDodge: this.getSuppressDodge(),
			parryHaste: this.getParryHaste(),
			spellSchool: this.getSpellSchool(),
			stats: this.stats.asArray(),
		});
	}

	fromProto(eventID: EventID, proto: TargetProto) {
		TypedEvent.freezeAllAndDo(() => {
			this.setId(eventID, proto.id);
			this.setName(eventID, proto.name);
			this.setLevel(eventID, proto.level);
			this.setMobType(eventID, proto.mobType);
			this.setTankIndex(eventID, proto.tankIndex);
			this.setSwingSpeed(eventID, proto.swingSpeed);
			this.setMinBaseDamage(eventID, proto.minBaseDamage);
			this.setDualWield(eventID, proto.dualWield);
			this.setDualWieldPenalty(eventID, proto.dualWieldPenalty);
			this.setCanCrush(eventID, proto.canCrush);
			this.setSuppressDodge(eventID, proto.suppressDodge);
			this.setParryHaste(eventID, proto.parryHaste);
			this.setSpellSchool(eventID, proto.spellSchool);
			this.setStats(eventID, new Stats(proto.stats));
		});
	}

	clone(eventID: EventID): Target {
		const newTarget = new Target(this.sim);
		newTarget.fromProto(eventID, this.toProto());
		return newTarget;
	}

	static defaultProto(): TargetProto {
		return TargetProto.create({
			level: Mechanics.BOSS_LEVEL,
			mobType: MobType.MobTypeDemon,
			tankIndex: 0,
			swingSpeed: 2,
			minBaseDamage: 10000,
			dualWield: false,
			dualWieldPenalty: false,
			canCrush: true,
			suppressDodge: false,
			parryHaste: true,
			spellSchool: SpellSchool.SpellSchoolPhysical,
			stats: Stats.fromMap({
				[Stat.StatArmor]: 7683,
				[Stat.StatBlockValue]: 54,
				[Stat.StatAttackPower]: 320,
			}).asArray(),
		});
	}

	static fromDefaults(eventID: EventID, sim: Sim): Target {
		const target = new Target(sim);
		target.fromProto(eventID, Target.defaultProto());
		return target;
	}
}
