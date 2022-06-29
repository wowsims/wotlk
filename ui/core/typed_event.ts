// An event ID uniquely identifies a single event that occurred, usually due to
// some user action like changing a piece of gear.
//
// Event IDs allow us to make sure that hierarchies of TypedEvents fire only once,
// for a given event. This is very important for certain features, like undo/redo.
export type EventID = number;

export interface Disposable {
	dispose(): void;
}

export interface Listener<T> {
	(eventID: EventID, event: T): any;
}

interface FiredEventData {
	eventID: EventID,
	error: Error,
}

interface FrozenEventData<T> {
	eventID: EventID,
	event: T,
}

/** Provides a type-safe event interface. */
export class TypedEvent<T> {
	// Optional label to help debug.
	private label: string;

	constructor(label?: string) {
		this.label = label || '';
	}

	private listeners: Array<Listener<T>> = [];

	// The events which have already been fired from this TypedEvent.
	private firedEvents: Array<FiredEventData> = [];

	// Currently frozen events pending on this TypedEvent. See freezeAll()
	// for more details.
	private frozenEvents: Array<FrozenEventData<T>> = [];

	// Registers a new listener to this event.
	on(listener: Listener<T>): Disposable {
		this.listeners.push(listener);
		return {
			dispose: () => this.off(listener),
		};
	}

	// Removes a listener from this event.
	off(listener: Listener<T>) {
		const idx = this.listeners.indexOf(listener);
		if (idx != -1) {
			this.listeners.splice(idx, 1);
		}
	}

	// Convenience for on() which calls off() autmatically after firing once.
	once(listener: Listener<T>): Disposable {
		const onceListener = (eventID: EventID, event: T) => {
			this.off(onceListener);
			listener(eventID, event);
		};

		return this.on(onceListener);
	}

	emit(eventID: EventID, event: T) {
		const originalEvent = this.firedEvents.find(fe => fe.eventID == eventID);
		if (originalEvent) {
			if (!thawing) {
				// Uncomment this for debugging TypedEvent stuff. There are a few legitimate
				// cases where it fires though and it can be very noisy.
				//console.warn('EventID collision outside of thawing, original event: ' + (originalEvent.error.stack || originalEvent.error));
			}
			return;
		}
		this.firedEvents.push({
			eventID: eventID,
			error: new Error('Original event'),
		});

		if (freezeCount > 0) {
			if (this.frozenEvents.length == 0) {
				frozenTypedEvents.push(this);
			}
			this.frozenEvents.push({
				eventID: eventID,
				event: event,
			});
		} else {
			this.fireEventInternal(eventID, event);
		}
	}

	private fireEventInternal(eventID: EventID, event: T) {
		this.listeners.forEach(listener => listener(eventID, event));
	}

	// Executes the provided callback while all TypedEvents are frozen.
	// Freezes all TypedEvent objects so that new calls to emit() do not fire the event.
	// Instead, the events will be held until the execution is finishd, at which point
	// all TypedEvents will fire all of the events that were frozen.
	//
	// This is used when a single user action activates multiple separate events, to ensure
	// none of them fire until all changes have been applied.
	//
	// This function is very similar to a locking mechanism.
	static freezeAllAndDo(func: () => void) {
		freezeCount++;

		try {
			func();
		} catch (e) {
			console.error('Caught error in freezeAllAndDo: ' + e);
		} finally {
			freezeCount--;
			if (freezeCount > 0) {
				// Don't do anything until things are fully unfrozen.
				return;
			}

			thawing = true;
			const typedEvents = frozenTypedEvents.slice();
			frozenTypedEvents = [];

			typedEvents.forEach(typedEvent => {
				const frozenEvents = typedEvent.frozenEvents.slice();
				typedEvent.frozenEvents = [];

				frozenEvents.forEach(frozenEvent => typedEvent.fireEventInternal(frozenEvent.eventID, frozenEvent.event));
			});
			thawing = false;
		}
	}

	static nextEventID(): EventID {
		return nextEventID++;
	}

	static onAny(events: Array<TypedEvent<any>>, label?: string): TypedEvent<void> {
		const newEvent = new TypedEvent<void>(label);
		events.forEach(emitter => emitter.on(eventID => newEvent.emit(eventID)));
		return newEvent;
	}
}

// If this is > 0 then events are frozen.
let freezeCount = 0;

// Indicates whether we are currently in the process of unfreezing. Just used to add a warning.
let thawing = false;

let frozenTypedEvents: Array<TypedEvent<any>> = [];
let nextEventID: EventID = 0;
