import * as Popper from '@popperjs/core';
import * as bootstrap from 'bootstrap';
import tippy from 'tippy.js';

import './shared/bootstrap_overrides';

declare global {
	interface Window {
		Popper: any;
		bootstrap: any;
		tippy: any;
	}
}

window.Popper = Popper;
window.bootstrap = bootstrap;
window.tippy = tippy;

function docReady(fn: any) {
	// see if DOM is already available
	if (document.readyState === "complete" || document.readyState === "interactive") {
		// call on next available tick
		setTimeout(fn, 1);
		// initializeCustomDropdownHandling();
	} else {
		document.addEventListener("DOMContentLoaded", fn);
	}
}

docReady(function() {
	document.getElementsByTagName('body')[0].classList.add('ready');
});
