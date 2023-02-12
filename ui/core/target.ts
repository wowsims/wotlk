import { MobType, TargetInput } from './proto/common.js';
import { SpellSchool } from './proto/common.js';
import { Stat } from './proto/common.js';
import { Target as TargetProto } from './proto/common.js';
import { PresetTarget } from './proto/common.js';
import { Stats } from './proto_utils/stats.js';

import * as Mechanics from './constants/mechanics.js';

import { Sim } from './sim.js';
import { EventID, TypedEvent } from './typed_event.js';
import { TargetInputs } from './target_inputs.js';

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
	private suppressDodge: boolean = false;
	private parryHaste: boolean = true;
	private tightEnemyDamage: boolean = false;
	private spellSchool: SpellSchool = SpellSchool.SpellSchoolPhysical;
	private targetInputs: TargetInputs = new TargetInputs()

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

	getTightEnemyDamage(): boolean {
		return this.tightEnemyDamage;
	}

	setTightEnemyDamage(eventID: EventID, newTightEnemyDamage: boolean) {
		if (newTightEnemyDamage == this.tightEnemyDamage)
			return;

		this.tightEnemyDamage = newTightEnemyDamage;
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

	getTargetInputs(): TargetInputs {
		return this.targetInputs;
	}

	setTargetInputs(eventID: EventID, newTargetInputs?: TargetInputs) {
		if (newTargetInputs?.equals(this.targetInputs))
			return;

		this.targetInputs = newTargetInputs ?? new TargetInputs();
		this.propChangeEmitter.emit(eventID);
	}

	hasTargetInputs(): boolean {
		return this.targetInputs.hasInputs();
	}

	getTargetInputsLength(): number {
		return this.targetInputs.getLength();
	}

	getTargetInputNumberValue(index: number): number {
		return this.targetInputs.getTargetInput(index)?.numberValue;
	}

	setTargetInputNumberValue(eventID: EventID, index: number, newValue: number) {
		if (this.getTargetInputNumberValue(index) == newValue)
			return;

		this.targetInputs.getTargetInput(index).numberValue = newValue;
		this.propChangeEmitter.emit(eventID);
	}

	getTargetInputBooleanValue(index: number): boolean {
		return this.targetInputs.getTargetInput(index)?.boolValue;
	}

	setTargetInputBooleanValue(eventID: EventID, index: number, newValue: boolean) {
		if (this.getTargetInputBooleanValue(index) == newValue)
			return;

		this.targetInputs.getTargetInput(index).boolValue = newValue;
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
			suppressDodge: this.getSuppressDodge(),
			parryHaste: this.getParryHaste(),
			tightEnemyDamage: this.getTightEnemyDamage(),
			spellSchool: this.getSpellSchool(),
			stats: this.stats.asArray(),
			targetInputs: this.targetInputs.asArray(),
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
			this.setSuppressDodge(eventID, proto.suppressDodge);
			this.setParryHaste(eventID, proto.parryHaste);
			this.setTightEnemyDamage(eventID, proto.tightEnemyDamage);
			this.setSpellSchool(eventID, proto.spellSchool);
			this.setTargetInputs(eventID, new TargetInputs(proto.targetInputs));
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
				[Stat.StatAttackPower]: 574,
			}).asArray(),
			targetInputs: new Array<TargetInput>(0),
		});
	}

	static fromDefaults(eventID: EventID, sim: Sim): Target {
		const target = new Target(sim);
		target.fromProto(eventID, Target.defaultProto());
		return target;
	}
}
