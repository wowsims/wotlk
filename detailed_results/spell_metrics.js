import { ActionMetrics } from '/wotlk/core/proto_utils/sim_result.js';
import { bucket } from '/wotlk/core/utils.js';
import { ColumnSortType, MetricsTable } from './metrics_table.js';
export class SpellMetricsTable extends MetricsTable {
    constructor(config) {
        config.rootCssClass = 'spell-metrics-root';
        super(config, [
            MetricsTable.nameCellConfig((metric) => {
                return {
                    name: metric.name,
                    actionId: metric.actionId,
                };
            }),
            {
                name: 'DPS',
                tooltip: 'Damage / Encounter Duration',
                sort: ColumnSortType.Descending,
                getValue: (metric) => metric.dps,
                getDisplayString: (metric) => metric.dps.toFixed(1),
            },
            {
                name: 'Avg Cast',
                tooltip: 'Damage / Casts',
                getValue: (metric) => metric.avgCast,
                getDisplayString: (metric) => metric.avgCast.toFixed(1),
            },
            {
                name: 'Avg Hit',
                tooltip: 'Damage / Hits',
                getValue: (metric) => metric.avgHit,
                getDisplayString: (metric) => metric.avgHit.toFixed(1),
            },
            {
                name: 'TPS',
                tooltip: 'Threat / Encounter Duration',
                columnClass: 'threat-metrics',
                getValue: (metric) => metric.tps,
                getDisplayString: (metric) => metric.tps.toFixed(1),
            },
            {
                name: 'Avg Cast',
                tooltip: 'Threat / Casts',
                columnClass: 'threat-metrics',
                getValue: (metric) => metric.avgCastThreat,
                getDisplayString: (metric) => metric.avgCastThreat.toFixed(1),
            },
            {
                name: 'Avg Hit',
                tooltip: 'Threat / Hits',
                columnClass: 'threat-metrics',
                getValue: (metric) => metric.avgHitThreat,
                getDisplayString: (metric) => metric.avgHitThreat.toFixed(1),
            },
            {
                name: 'Casts',
                tooltip: 'Casts',
                getValue: (metric) => metric.casts,
                getDisplayString: (metric) => metric.casts.toFixed(1),
            },
            {
                name: 'Hits',
                tooltip: 'Hits',
                getValue: (metric) => metric.landedHits,
                getDisplayString: (metric) => metric.landedHits.toFixed(1),
            },
            {
                name: 'Miss %',
                tooltip: 'Misses / (Hits + Misses)',
                getValue: (metric) => metric.missPercent,
                getDisplayString: (metric) => metric.missPercent.toFixed(2) + '%',
            },
            {
                name: 'Crit %',
                tooltip: 'Crits / Hits',
                getValue: (metric) => metric.critPercent,
                getDisplayString: (metric) => metric.critPercent.toFixed(2) + '%',
            },
        ]);
    }
    customizeRowElem(action, rowElem) {
        if (action.hitAttempts == 0 && action.dps == 0) {
            rowElem.classList.add('threat-metrics');
        }
    }
    getGroupedMetrics(resultData) {
        const players = resultData.result.getPlayers(resultData.filter);
        if (players.length != 1) {
            return [];
        }
        const player = players[0];
        const actions = player.getSpellActions().map(action => action.forTarget(resultData.filter));
        const actionGroups = ActionMetrics.groupById(actions);
        const petsByName = bucket(player.pets, pet => pet.name);
        const petGroups = Object.values(petsByName).map(pets => ActionMetrics.joinById(pets.map(pet => pet.getSpellActions().map(action => action.forTarget(resultData.filter))).flat(), true));
        return actionGroups.concat(petGroups);
    }
    mergeMetrics(metrics) {
        return ActionMetrics.merge(metrics, true, metrics[0].unit?.petActionId || undefined);
    }
    shouldCollapse(metric) {
        return !metric.unit?.isPet;
    }
}
