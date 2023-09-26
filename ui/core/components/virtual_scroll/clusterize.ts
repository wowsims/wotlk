// Based on clusterize
/* Clusterize.js - v1.0.0 - 2023-01-22
 http://NeXTs.github.com/Clusterize.js/
 Copyright (c) 2015 Denis Lukov; Licensed MIT */

interface ClusterizeCallbacks {
    clusterWillChange?: () => void;
    clusterChanged?: () => void;
    scrollingProgress?: (n:number) => void;
}

interface ClusterizeParams {
    rows: Element[];
    scroll_elem: HTMLElement;
    content_elem: HTMLElement;

    rows_in_block?: number;
    blocks_in_cluster?: number;
    show_no_data_row?: boolean;
    no_data_class?:string;
    no_data_text?:string;
    keep_parity?:boolean;
    tag?:keyof HTMLElementTagNameMap;
    callbacks?: ClusterizeCallbacks;
    item_height?:number;
}

let defaults = {
    rows_in_block: 50,
    blocks_in_cluster: 4,
    tag: null,
    show_no_data_row: true,
    no_data_class: 'clusterize-no-data',
    no_data_text: 'No data',
    keep_parity: true,
    callbacks: {}
}

function getStyle(prop:string, elem:HTMLElement) {
    return window.getComputedStyle(elem)[prop as any]
}

interface ClusterizeOptions {
    rows_in_block: number;
    blocks_in_cluster: number;
    show_no_data_row: boolean;
    no_data_class:string;
    no_data_text:string;
    keep_parity:boolean;
    tag:keyof HTMLElementTagNameMap | null;
    callbacks: ClusterizeCallbacks;

    cluster_height: number;
    item_height: number;
    content_tag: string;

    block_height: number;
    scroll_top: number;
    rows_in_cluster: number;
}

interface ClusterizeDatasource {
    getNumberOfRows: () => number;
    generateRows: (startIdx: number, endIdx: number) => Element[];
}

export class Clusterize {
    private options: ClusterizeOptions;
    private scroll_elem: Element;
    private content_elem: HTMLElement;
    private displayedRows: Element[];

    private last_cluster = 0;
    private scroll_debounce = 0;
    private pointer_events_set = false;
    private resize_debounce = 0;

    private ds: ClusterizeDatasource;
    private cache = {};

    constructor(ds: ClusterizeDatasource, params: ClusterizeParams) {
        this.options = {
            rows_in_block: params.rows_in_block !== undefined ? params.rows_in_block : defaults.rows_in_block,
            blocks_in_cluster: params.blocks_in_cluster !== undefined ? params.blocks_in_cluster : defaults.blocks_in_cluster,
            show_no_data_row: params.show_no_data_row !== undefined ? params.show_no_data_row : defaults.show_no_data_row,
            no_data_class: params.no_data_class !== undefined ? params.no_data_class : defaults.no_data_class,
            no_data_text: params.no_data_text !== undefined ? params.no_data_text : defaults.no_data_text,
            keep_parity: params.keep_parity !== undefined ? params.keep_parity : defaults.keep_parity,
            tag: params.tag !== undefined ? params.tag : defaults.tag,
            callbacks: params.callbacks !== undefined ? params.callbacks : defaults.callbacks,
            cluster_height: 0,
            item_height: params.item_height !== undefined ? params.item_height : 0,
            content_tag: '',
            block_height: 0,
            scroll_top: 0,
            rows_in_cluster: 0,
        }

        this.ds = ds;
        this.displayedRows = params.rows;
        this.scroll_elem = params.scroll_elem;
        this.content_elem = params.content_elem;

        // tabindex forces the browser to keep focus on the scrolling list, fixes #11
        if(!this.content_elem.hasAttribute('tabindex'))
            this.content_elem.setAttribute('tabindex', '0');
  
        let scroll_top = this.scroll_elem.scrollTop;
  
        // append initial data
        this.getRowsHeight(params.rows);
        this.insertToDom(params.rows, this.cache);
  
        // restore the scroll position
        this.scroll_elem.scrollTop = scroll_top;
    
        // adding scroll handler
        this.scroll_elem.addEventListener('scroll', this.scrollEv.bind(this));
        window.addEventListener('resize', this.resizeEv.bind(this));
    }

    destroy(clean? : boolean) {
        this.scroll_elem.removeEventListener('scroll', this.scrollEv.bind(this));
        window.removeEventListener('resize', this.resizeEv.bind(this));
        if (clean)
            this.setContentElemRows(this.generateEmptyRow());
    }

    dispose() {
        this.destroy();
    }

    refresh(force?: boolean) {
        if(this.getRowsHeight(this.displayedRows) || force)
            this.update();
    }

    update() {
        let scroll_top = this.scroll_elem.scrollTop;
        if(this.ds.getNumberOfRows() * this.options.item_height < scroll_top) {
            this.scroll_elem.scrollTop = 0;
            this.last_cluster = 0;
        }
        this.insertToDom(this.displayedRows, this.cache);
        this.scroll_elem.scrollTop = scroll_top;
    }

    elementUpdate(cb: (e: Element) => void) {
        for (let r of this.displayedRows) {
            cb(r);
        }
    }

    getScrollProgress() {
        return this.options.scroll_top / (this.ds.getNumberOfRows() * this.options.item_height) * 100 || 0;
    }

    private scrollEv() {
        // fixes scrolling issue on Mac #3
        let is_mac = navigator.platform.toLowerCase().indexOf('mac') + 1;
        if (is_mac) {
            if (!this.pointer_events_set) 
                this.content_elem.style.pointerEvents = 'none';
            this.pointer_events_set = true;
            clearTimeout(this.scroll_debounce);
            this.scroll_debounce = window.setTimeout(() => {
                this.content_elem.style.pointerEvents = 'auto';
                this.pointer_events_set = false;
            }, 50);
        }
        if (this.last_cluster != (this.last_cluster = this.getClusterNum(this.ds.getNumberOfRows())))
            this.insertToDom(this.displayedRows, this.cache);
        if (this.options.callbacks.scrollingProgress)
            this.options.callbacks.scrollingProgress(this.getScrollProgress());
    };

    private resizeEv() {
        clearTimeout(this.resize_debounce);
        this.resize_debounce = window.setTimeout(this.refresh.bind(this), 100);
    }

    private insertToDom(rows: Element[], cache: any) {
        // explore row's height
        if(!this.options.cluster_height) {
            this.exploreEnvironment(rows, cache);
        }
        let data = this.generate();
        let this_cluster_rows = data.rows;
        this.displayedRows = data.rows;
        let this_cluster_content_changed = this.checkChanges('data', this_cluster_rows, cache);
        let top_offset_changed = this.checkChanges('top', data.top_offset, cache);
        let only_bottom_offset_changed = this.checkChanges('bottom', data.bottom_offset, cache);
        let callbacks = this.options.callbacks;
        let layout: Element[] = [];

        if(this_cluster_content_changed || top_offset_changed) {
            if(data.top_offset) {
                if (this.options.keep_parity) {
                    let parity = this.renderExtraTag('keep-parity');
                    parity.hidden = true;
                    layout.push(parity);
                }
                layout.push(this.renderExtraTag('top-space', data.top_offset));
            }
            layout.push(...this_cluster_rows);
            data.bottom_offset && layout.push(this.renderExtraTag('bottom-space', data.bottom_offset));
            callbacks.clusterWillChange && callbacks.clusterWillChange();
            this.setContentElemRows(layout);
            this.options.content_tag == 'ol' && this.content_elem.setAttribute('start', data.rows_above.toString());
            this.content_elem.style.counterIncrement = 'clusterize-counter ' + (data.rows_above-1);
            callbacks.clusterChanged && callbacks.clusterChanged();
        } else if(only_bottom_offset_changed) {
            (this.content_elem.lastChild as HTMLElement).style.height = data.bottom_offset + 'px';
        }
    }

    // get tag name, content tag name, tag height, calc cluster height
    private exploreEnvironment(rows: Element[], cache: any) {
        let opts = this.options;
        opts.content_tag = this.content_elem.tagName.toLowerCase();
        if (!rows.length)
            return;
        if (this.content_elem.children.length <= 1)
            this.setContentElemRows([rows[0], rows[0], rows[0]]);
        if (!opts.tag)
            opts.tag = this.content_elem.children[0].tagName.toLowerCase() as keyof HTMLElementTagNameMap;
            this.getRowsHeight(rows);
    }

    private setContentElemRows(rows: Element[]) {
        this.content_elem.innerHTML = ``;
        this.content_elem.append(...rows);
    }

    private getRowsHeight(rows: Element[]) {
        let opts = this.options;
        let prev_item_height = opts.item_height;
        opts.cluster_height = 0;
        let nodes = this.content_elem.children;
        if (rows.length && nodes.length) {
            let node = nodes[Math.floor(nodes.length / 2)] as HTMLElement;
            opts.item_height = node.offsetHeight;
            // consider table's border-spacing
            if(opts.tag == 'tr' && getStyle('borderCollapse', this.content_elem) != 'collapse')
                opts.item_height += parseInt(getStyle('borderSpacing', this.content_elem), 10) || 0;
            // consider margins (and margins collapsing)
            if(opts.tag != 'tr') {
                let marginTop = parseInt(getStyle('marginTop', node), 10) || 0;
                let marginBottom = parseInt(getStyle('marginBottom', node), 10) || 0;
                opts.item_height += Math.max(marginTop, marginBottom);
            }
        }
        opts.block_height = opts.item_height * opts.rows_in_block;
        opts.rows_in_cluster = opts.blocks_in_cluster * opts.rows_in_block;
        opts.cluster_height = opts.blocks_in_cluster * opts.block_height;
        return prev_item_height != opts.item_height;
    }

    // generate cluster for current scroll position
    private generate() {
        let opts = this.options;
        let rows_len = this.ds.getNumberOfRows();
        if (rows_len < opts.rows_in_block) {
            let rows = this.ds.generateRows(0, rows_len);
            return {
                top_offset: 0,
                bottom_offset: 0,
                rows_above: 0,
                rows: rows_len ? rows : this.generateEmptyRow()
            }
        }
        let items_start = Math.max((opts.rows_in_cluster - opts.rows_in_block) * this.getClusterNum(rows_len), 0),
          items_end = items_start + opts.rows_in_cluster,
          top_offset = Math.max(items_start * opts.item_height, 0),
          bottom_offset = Math.max((rows_len - items_end) * opts.item_height, 0),
          this_cluster_rows = [],
          rows_above = items_start;
        if(top_offset < 1) {
            rows_above++;
        }
        this_cluster_rows = this.ds.generateRows(items_start, items_end);
        return {
            top_offset: top_offset,
            bottom_offset: bottom_offset,
            rows_above: rows_above,
            rows: this_cluster_rows
        }
    }

    // generate empty row if no data provided
    private generateEmptyRow() : Element[] {
        let opts = this.options;
        if(!opts.tag || !opts.show_no_data_row) 
            return [];
        let empty_row = document.createElement(opts.tag);
        let no_data_content = document.createTextNode(opts.no_data_text), td;
        empty_row.className = opts.no_data_class;
        if(opts.tag == 'tr') {
            td = document.createElement('td');
            // fixes #53
            td.colSpan = 100;
            td.appendChild(no_data_content);
        }
        empty_row.appendChild(td || no_data_content);
        return [empty_row];
    }

    // get current cluster number
    private getClusterNum(num_rows: number) {
        let opts = this.options;
        opts.scroll_top = this.scroll_elem.scrollTop;
        let cluster_divider = Math.max(opts.cluster_height - opts.block_height, 1);
        let current_cluster = Math.floor(opts.scroll_top / cluster_divider);
        let max_cluster = Math.floor((num_rows * opts.item_height) / cluster_divider);
        return Math.min(current_cluster, max_cluster);
    }

    private renderExtraTag(class_name: string, height?: number) {
        let tag = document.createElement(this.options.tag!);
        let clusterize_prefix = 'clusterize-';

        tag.className = [clusterize_prefix + 'extra-row', clusterize_prefix + class_name].join(' ');
        tag.style.marginTop = '0';
        tag.style.marginBottom = '0';
        height && (tag.style.height = height + 'px');
        return tag;
    }

    private checkChanges(type: string, value: any, cache: any) {
        var changed = value != cache[type];
        cache[type] = value;
        return changed;
    }
}