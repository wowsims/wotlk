import * as bootstrap from 'bootstrap';

document.querySelectorAll('[data-bs-toggle="dropdown"]').forEach(toggle => {
  let dropdown = new bootstrap.Dropdown(toggle);
  let dropdownMenu = toggle.nextElementSibling as HTMLElement;

  toggle.addEventListener('mouseover', event => {
    dropdown.show();
  });

  toggle.addEventListener('mouseleave', event => {
    let e = event as MouseEvent;
    let target = e.relatedTarget as HTMLElement;

    if (target != dropdownMenu && !target.closest('.dropdown-menu'))
      dropdown.hide();
  });

  dropdownMenu.addEventListener('mouseleave', event => {
    let e = event as MouseEvent;
    let target = e.relatedTarget as HTMLElement;

    if (target != toggle)
      dropdown.hide();
  });
});
