import { Player } from "../../player";
import { Spec, Stat } from "../../proto/common";

import { IconPicker, IconPickerConfig } from "../icon_picker";
import { MultiIconPicker, MultiIconPickerConfig } from "../multi_icon_picker";
import { IconEnumPicker, IconEnumPickerConfig } from "../icon_enum_picker";
import { IndividualSimUI } from "ui/core/individual_sim_ui";

export interface StatOption<PickerType, ConfigType> {
	config: ConfigType,
	stats: Array<Stat>,
	picker?: PickerType,
}

export interface IconPickerStatOption extends StatOption<
	typeof IconPicker<Player<any>, any>,
	IconPickerConfig<Player<any>, any>
> {}

export interface MultiIconPickerStatOption extends StatOption<
	typeof MultiIconPicker<Player<any>>,
	MultiIconPickerConfig<Player<any>>
> {}

export interface IconEnumPickerStatOption extends StatOption<
  typeof IconEnumPicker<Player<any>, any>,
  IconEnumPickerConfig<Player<any>, any>
> {}

export type StatOptions = Array<IconPickerStatOption | MultiIconPickerStatOption | IconEnumPickerStatOption>

export function relevantStatOptions(options: StatOptions, simUI: IndividualSimUI<Spec>): StatOptions {
  return options
    .filter(option =>
      option!.stats.length == 0 ||
      option!.stats.some(stat => simUI.individualConfig.epStats.includes(stat)))
}
