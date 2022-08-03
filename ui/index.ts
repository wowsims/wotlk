import { NotificationDOMFactory } from "./core/notifications/dom_factory";
import { Notifier } from "./core/notifications/notifier";

let notifier: Notifier;

window.onload = () => {
	const notification_wrapper = document.createElement("div");
	notification_wrapper.classList.add("notifications");
	document.body.appendChild(notification_wrapper);
	const notification_dom_factory = new NotificationDOMFactory(document);
	if (notification_wrapper !== null) {
		notifier = new Notifier(notification_wrapper, notification_dom_factory);
	} else {
		console.error("Element not found: #notifications")
	}

	setTimeout(() => {
		notifier.notify_debug("Hi");
	}, 1000)
	setTimeout(() => {
		notifier.notify_info("Hi");
	}, 3000)
	setTimeout(() => {
		notifier.notify_warning("Hi");
	}, 4000)
	setTimeout(() => {
		notifier.notify_error("Hi");
	}, 5000)
}
console.log("INIT")
