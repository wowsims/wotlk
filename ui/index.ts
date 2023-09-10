import * as Popper from '@popperjs/core';
import * as bootstrap from 'bootstrap';
import tippy from 'tippy.js';

import './shared/bootstrap_overrides';
import './shared/gear_compare';

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

// Force scroll to top when refreshing
if (history.scrollRestoration) {
	history.scrollRestoration = 'manual';
} else {
	window.onbeforeunload = function () {
		window.scrollTo(0, 0);
	}
}

function docReady(fn: any) {
	// see if DOM is already available
	if (document.readyState === "complete" || document.readyState === "interactive") {
		// call on next available tick
		setTimeout(fn, 1);
	} else {
		document.addEventListener("DOMContentLoaded", fn);
	}
}

docReady(function() {
	document.body.classList.add('ready');
});
