import { Spec } from './proto/common.js';

// This file is for anything related to launching a new sim. DO NOT touch this
// file until your sim is ready to launch!

export enum LaunchStatus {
    Unlaunched,
    Alpha,
    Beta,
    Launched,
}

export const raidSimLaunched = false;

// This list controls which links are shown in the top-left dropdown menu.
export const simLaunchStatuses: Record<Spec, LaunchStatus> = {
    [Spec.SpecBalanceDruid]: LaunchStatus.Unlaunched,
    [Spec.SpecElementalShaman]: LaunchStatus.Unlaunched,
    [Spec.SpecEnhancementShaman]: LaunchStatus.Unlaunched,
    [Spec.SpecFeralDruid]: LaunchStatus.Unlaunched,
    [Spec.SpecFeralTankDruid]: LaunchStatus.Unlaunched,
    [Spec.SpecHunter]: LaunchStatus.Alpha,
    [Spec.SpecMage]: LaunchStatus.Unlaunched,
    [Spec.SpecRogue]: LaunchStatus.Unlaunched,
    [Spec.SpecRetributionPaladin]: LaunchStatus.Unlaunched,
    [Spec.SpecProtectionPaladin]: LaunchStatus.Unlaunched,
    [Spec.SpecShadowPriest]: LaunchStatus.Unlaunched,
    [Spec.SpecWarlock]: LaunchStatus.Unlaunched,
    [Spec.SpecWarrior]: LaunchStatus.Unlaunched,
    [Spec.SpecProtectionWarrior]: LaunchStatus.Unlaunched,
    [Spec.SpecSmitePriest]: LaunchStatus.Unlaunched,
    [Spec.SpecDeathknight]: LaunchStatus.Unlaunched,
    [Spec.SpecTankDeathknight]: LaunchStatus.Unlaunched,
};

export function getLaunchedSims(): Array<Spec> {
    return Object.keys(simLaunchStatuses)
        .map(specStr => parseInt(specStr) as Spec)
        .filter(spec => simLaunchStatuses[spec] > LaunchStatus.Unlaunched);
}
