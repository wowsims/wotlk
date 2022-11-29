import * as bootstrap from 'bootstrap';
import { isDescendant } from './utils';

let body = document.querySelector('body') as HTMLElement;

// Custom dropdown event handlers for mouseover dropdowns
body.addEventListener('mouseover', event => {
  let target = event.target as HTMLElement;
  let toggle = target.closest('[data-bs-toggle=dropdown]');
  if (toggle) {
    let dropdown = bootstrap.Dropdown.getOrCreateInstance(toggle);
    dropdown.show();
  }
}, true);

body.addEventListener('mouseleave', event => {
  let e = event as MouseEvent;
  let target = event.target as HTMLElement;
  let toggle = target.closest('[data-bs-toggle=dropdown]');
  if (toggle) {
    let dropdown = bootstrap.Dropdown.getOrCreateInstance(toggle);
    let dropdownMenu = toggle.nextElementSibling as HTMLElement;
    let relatedTarget = e.relatedTarget as HTMLElement;
    if (!isDescendant(relatedTarget, dropdownMenu))
      dropdown.hide();
  }
}, true);

body.addEventListener('mouseleave', event => {
  let e = event as MouseEvent;
  let target = event.target as HTMLElement;
  let dropdownMenu = target.closest('.dropdown-menu') as HTMLElement;
  if (dropdownMenu) {
    let toggle = dropdownMenu.previousElementSibling as HTMLElement;
    let dropdown = bootstrap.Dropdown.getOrCreateInstance(toggle);
    let relatedTarget = e.relatedTarget as HTMLElement;
    if (!isDescendant(relatedTarget, dropdownMenu) && e.relatedTarget != toggle)
      dropdown.hide();
  }
}, true);
