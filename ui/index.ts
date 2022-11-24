import * as Popper from '@popperjs/core';
import * as bootstrap from 'bootstrap';

declare global {
  interface Window {
    Popper: any;
    bootstrap: any;
  }
}

window.Popper = Popper;
window.bootstrap = bootstrap;

import './shared/bootstrap_overrides';
