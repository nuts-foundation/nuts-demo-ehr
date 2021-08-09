# eOverdracht FHIR Specification Gaps

* What HTTP response code should be returned by the notification endpoint?
* What if the notification endpoint returns non-OK HTTP status code (< 199 || > 299)
* What content-type should the request have? (none?)
* Should the notification be retried?
* eOverdracht Bolt states:
  * Het endpoint waar de notificatie heen moet is een combinatie van het base endpoint en het relatieve pad zoals gedefinieerd in het TO van Nictiz.
  * The relative path isn't defined in the Nictiz spec

# eOverdracht Bolt Specification Gaps
* The Bolt states: "Omdat het security token geen gebruikersinformatie bevat, mogen er nooit persoonsgegevens meegestuurd worden in de Task."
  * But this is allowed by the Nictiz specification (but optional)?