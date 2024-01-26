import { IndividualSimUI } from "../../individual_sim_ui";
import { Player } from "../../player";
import { Faction, Spec, Stat } from "../../proto/common";
import { ActionId } from "../../proto_utils/action_id";

import { IconEnumPicker, IconEnumPickerConfig } from "../icon_enum_picker";
import { IconPicker, IconPickerConfig } from "../icon_picker";
import { MultiIconPicker, MultiIconPickerConfig } from "../multi_icon_picker";

export interface ItemInputConfig {
	id: ActionId
	minLevel?: number
	maxLevel?: number
	faction?: Faction
}

export interface StatOption {
	stats: Array<Stat>,
}

export interface ItemStatOption extends StatOption {
	item: ItemInputConfig,
}

export interface PickerStatOption<PickerType, ConfigType> extends StatOption {
	config: ConfigType,
	picker: PickerType,
}

export interface IconPickerStatOption extends PickerStatOption<
	typeof IconPicker<Player<any>, any>,
	IconPickerConfig<Player<any>, any>
> {}

export interface MultiIconPickerStatOption extends PickerStatOption<
	typeof MultiIconPicker<Player<any>>,
	MultiIconPickerConfig<Player<any>>
> {}

export interface IconEnumPickerStatOption extends PickerStatOption<
  typeof IconEnumPicker<Player<any>, any>,
  IconEnumPickerConfig<Player<any>, any>
> {}

export type ItemStatOptions = ItemStatOption
export type PickerStatOptions = IconPickerStatOption | MultiIconPickerStatOption | IconEnumPickerStatOption
export type StatOptions<Options extends ItemStatOptions | PickerStatOptions> = Array<Options>

export function relevantStatOptions<OptionsType extends ItemStatOptions | PickerStatOptions>(
	options: StatOptions<OptionsType>,
	simUI: IndividualSimUI<Spec>
): StatOptions<OptionsType> {
  return options
    .filter(option =>
      option.stats.length == 0 ||
      option.stats.some(stat => simUI.individualConfig.epStats.includes(stat)))
}
