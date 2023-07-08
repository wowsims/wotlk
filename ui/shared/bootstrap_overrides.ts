import { Dropdown, Popover, Tooltip } from 'bootstrap';
import { isDescendant } from './utils';

Dropdown.Default.offset = [0,0];
//Dropdown.Default.display = "static";

Tooltip.Default.trigger = "hover";

let body = document.querySelector('body') as HTMLElement;

// Custom dropdown event handlers for mouseover dropdowns
body.addEventListener('mouseover', event => {
  let target = event.target as HTMLElement;
  let toggle = target.closest('[data-bs-toggle=dropdown]');
  if (toggle && !toggle.classList.contains('open-on-click')) {
    let dropdown = Dropdown.getOrCreateInstance(toggle);
    dropdown.show();
  }
}, true);

body.addEventListener('mouseleave', event => {
  let e = event as MouseEvent;
  let target = event.target as HTMLElement;
  let toggle = target.closest('[data-bs-toggle=dropdown]') as HTMLElement | null;
  // Hide dropdowns when hovering off of the toggle, so long as the new target is not part of the dropdown as well
  if (toggle) {
    let dropdown = Dropdown.getOrCreateInstance(toggle);
    let dropdownMenu = toggle.nextElementSibling as HTMLElement;
    let relatedTarget = e.relatedTarget as HTMLElement;
    if (relatedTarget == null || (!isDescendant(relatedTarget, dropdownMenu) && !isDescendant(relatedTarget, toggle)))
      dropdown.hide();
  }

  let dropdownMenu = target.closest('.dropdown-menu') as HTMLElement;
  // Hide dropdowns when hovering off of the menu, so long as the new target is not part of the dropdown as well
  if (dropdownMenu) {
    let toggle = dropdownMenu.previousElementSibling as HTMLElement;
    let dropdown = Dropdown.getOrCreateInstance(toggle);
    let relatedTarget = e.relatedTarget as HTMLElement;
    if (relatedTarget == null || (!isDescendant(relatedTarget, dropdownMenu) && e.relatedTarget != toggle))
      dropdown.hide();
  }
}, true);

let closePopovers = () => {
  document.querySelectorAll('[data-bs-toggle="popover"][aria-describedby]').forEach(e => {
    let p = Popover.getOrCreateInstance(e);
    p.hide();
  });
}

body.addEventListener('show.bs.popover', (event) => {
  closePopovers();

  document.querySelectorAll('[data-bs-toggle="tooltip"][aria-describedby]').forEach(e => {
    let t = Tooltip.getOrCreateInstance(e);
    t.hide();
  });

  document.querySelectorAll('.tooltip').forEach(e => e.remove());
}, true);

body.addEventListener('show.bs.tooltip', (event) => {
  document.querySelectorAll('[data-bs-toggle="tooltip"][aria-describedby]').forEach(e => {
    let t = Tooltip.getOrCreateInstance(e);
    t.hide();
  });

  document.querySelectorAll('.tooltip').forEach(e => e.remove());
}, true);

document.onkeydown = (event) => {
  event = event || window.event;

  if (event.key == 'Escape') {
    closePopovers();
  }
}
