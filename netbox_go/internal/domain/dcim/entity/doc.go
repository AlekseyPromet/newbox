// Package entity содержит доменные сущности DCIM модуля NetBox.
//
// Сущности представляют основные бизнес-объекты предметной области:
//   - Sites: Region, SiteGroup, Site, Location
//   - Racks: RackType, RackRole, Rack, RackReservation
//   - Devices: Manufacturer, DeviceType, Device, ModuleType, Module
//   - Components: ConsolePort, ConsoleServerPort, PowerPort, PowerOutlet, Interface и др.
//   - Cables: Cable, CableTermination
//   - Power: PowerFeed, PowerPanel
//
// Каждая сущность реализует метод Validate() для проверки бизнес-правил
// и следует принципам SOLID (единственная ответственность).
package entity
