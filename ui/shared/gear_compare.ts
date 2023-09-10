declare var WH: any;

const cloneTooltip = (): HTMLElement => {
  const wowheadTooltip = document.querySelector('.wowhead-tooltip:not(#currently-equipped-tooltip)') as HTMLElement;
  const clone = wowheadTooltip.cloneNode(true) as HTMLElement;
  clone.id = 'currently-equipped-tooltip'
  clone.style.display = 'none';
  clone.style.visibility = 'hidden';
  wowheadTooltip.after(clone);
  return clone
}

export const setGearCompare = (elem: HTMLElement) => {
  WH.Tooltips.triggerTooltip(elem);
}

export const initializeGearCompare = () => {
  const wowheadTooltip = document.querySelector('.wowhead-tooltip:not(#currently-equipped-tooltip)') as HTMLElement;
  const cloneTooltip = (): HTMLElement => {
    const clone = wowheadTooltip.cloneNode(true) as HTMLElement;
    clone.id = 'currently-equipped-tooltip'
    clone.style.display = 'none';
    clone.style.visibility = 'hidden';
    wowheadTooltip.after(clone);
    return clone
  }
  let currentEquipTooltip: HTMLElement = document.querySelector('#currently-equipped-tooltip') ?? cloneTooltip();

  document.querySelectorAll<HTMLElement>('[data-gear-compare="true"]').forEach((elem) => {
    elem.addEventListener('mouseenter', (event: Event) => {
      currentEquipTooltip.setAttribute('data-gear-compare-active', 'true');

      if ((event as KeyboardEvent).shiftKey) {
        currentEquipTooltip.style.display = 'block';
        currentEquipTooltip.style.visibility = 'visible';
      }
    })
    elem.addEventListener('mousemove', () => {
      currentEquipTooltip.style.left = `calc(${wowheadTooltip.style.left} + ${wowheadTooltip.style.width})`;
      currentEquipTooltip.style.top = wowheadTooltip.style.top;
    })
    elem.addEventListener('mouseleave', () => {
      currentEquipTooltip.setAttribute('data-gear-compare-active', 'false');
      currentEquipTooltip.style.display = 'none';
			currentEquipTooltip.style.visibility = 'hidden';
    })
  })
}

document.addEventListener('keydown', (event: KeyboardEvent) => {
  if (event.key === 'Shift') {
    const wowheadTooltip = document.querySelector<HTMLElement>('.wowhead-tooltip:not(#currently-equipped-tooltip)');
    const currentEquipTooltip = document.querySelector<HTMLElement>('#currently-equipped-tooltip');

    if (wowheadTooltip && currentEquipTooltip && currentEquipTooltip.getAttribute('data-gear-compare-active') === 'true') {
      currentEquipTooltip.style.display = 'block';
      currentEquipTooltip.style.visibility = 'visible';
    }
  }
})

document.addEventListener('keyup', (event: KeyboardEvent) => {
  if (event.key === 'Shift') {
    const currentEquipTooltip = document.querySelector<HTMLElement>('#currently-equipped-tooltip');

    if (currentEquipTooltip) {
      currentEquipTooltip.style.display = 'none';
			currentEquipTooltip.style.visibility = 'hidden';
    }
  }
})
