import * as InputHelpers from '/wotlk/core/components/input_helpers.js';
// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const StartingRunicPower = InputHelpers.makeSpecOptionsNumberInput({
    fieldName: 'startingRunicPower',
    label: 'Starting Runic Power',
    labelTooltip: 'Initial RP at the start of each iteration.',
});
export const DeathKnightRotationConfig = {
    inputs: [],
};
