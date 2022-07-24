import { Spec } from '/wotlk/core/proto/common.js';
// This file is for anything related to launching a new sim. DO NOT touch this
// file until your sim is ready to launch!
export var LaunchStatus;
(function (LaunchStatus) {
    LaunchStatus[LaunchStatus["Unlaunched"] = 0] = "Unlaunched";
    LaunchStatus[LaunchStatus["Alpha"] = 1] = "Alpha";
    LaunchStatus[LaunchStatus["Beta"] = 2] = "Beta";
    LaunchStatus[LaunchStatus["Launched"] = 3] = "Launched";
})(LaunchStatus || (LaunchStatus = {}));
export const raidSimLaunched = false;
// This list controls which links are shown in the top-left dropdown menu.
export const simLaunchStatuses = {
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
export function getLaunchedSims() {
    return Object.keys(simLaunchStatuses)
        .map(specStr => parseInt(specStr))
        .filter(spec => simLaunchStatuses[spec] > LaunchStatus.Unlaunched);
}
