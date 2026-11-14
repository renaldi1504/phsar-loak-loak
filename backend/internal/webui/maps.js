'use strict';
// Peta untuk Seller Center & admin.
// GOOGLE_MAPS_API_KEY diisi → Google Maps + pencarian tempat (Places API New) + alamat otomatis (Geocoding).
// Kosong → cadangan OpenStreetMap (Leaflet), pin dipasang dengan klik.
const SIHANOUKVILLE = {lat: 10.6274, lng: 103.5223};
const SERVICE_BOUNDS = {north: 10.72, south: 10.55, east: 103.62, west: 103.44}; // kotak kasar, sama dengan seed service_areas

const Maps = {
  cfg: window.PHSAR_CONFIG || {},
  google: () => !!(window.PHSAR_CONFIG || {}).googleMapsKey,
  _loading: null,

  load() {
    if (!this.google()) return Promise.resolve(false);
    if (window.google?.maps?.importLibrary) return Promise.resolve(true);
    if (!this._loading) this._loading = new Promise((resolve, reject) => {
      window.__phsarMapsReady = () => resolve(true);
      const s = document.createElement('script');
      const q = new URLSearchParams({key: this.cfg.googleMapsKey, v: 'weekly', loading: 'async', callback: '__phsarMapsReady',
        language: ({km: 'km', zh: 'zh-CN'})[P.lang] || 'en', region: 'KH'});
      s.src = 'https://maps.googleapis.com/maps/api/js?' + q;
      s.async = true;
      s.onerror = () => { this._loading = null; reject(new Error('Google Maps tidak bisa dimuat')); };
      document.head.append(s);
    });
    return this._loading;
  },

  // picker: peta dengan satu pin yang bisa digeser + kotak cari tempat.
  // opts: {mapEl, searchEl, pin:{lat,lng}|null, onPin(pos), onPlace({lat,lng,address,name,placeId}), onAddress(address)}
  async picker(opts) {
    let ok = false;
    try { ok = await this.load(); } catch (e) { console.warn(e); }
    return ok ? this._googlePicker(opts) : this._leafletPicker(opts);
  },

  async _googlePicker({mapEl, searchEl, pin, onPin, onPlace, onAddress}) {
    const [{Map}, {AdvancedMarkerElement}, places, {Geocoder}] = await Promise.all([
      google.maps.importLibrary('maps'), google.maps.importLibrary('marker'),
      google.maps.importLibrary('places'), google.maps.importLibrary('geocoding')]);
    const map = new Map(mapEl, {center: pin || SIHANOUKVILLE, zoom: pin ? 18 : 13, mapId: this.cfg.googleMapId || 'DEMO_MAP_ID',
      mapTypeControl: true, streetViewControl: false, fullscreenControl: false, gestureHandling: 'greedy'});
    const geocoder = new Geocoder();
    let marker = null;
    const toLit = p => typeof p.lat === 'function' ? {lat: p.lat(), lng: p.lng()} : {lat: p.lat, lng: p.lng};
    const round = p => ({lat: Number(p.lat.toFixed(6)), lng: Number(p.lng.toFixed(6))});
    const reverse = async pos => {
      try {
        const {results} = await geocoder.geocode({location: pos});
        if (results?.[0]) onAddress?.(results[0].formatted_address, results[0].place_id);
      } catch (e) { console.warn('geocode', e); }
    };
    const place = (pos, {geocode = true} = {}) => {
      pos = round(pos);
      if (!marker) {
        marker = new AdvancedMarkerElement({map, position: pos, gmpDraggable: true, title: 'Pickup'});
        marker.addListener('dragend', () => { const p = round(toLit(marker.position)); onPin(p); reverse(p); });
      } else marker.position = pos;
      onPin(pos);
      if (geocode) reverse(pos);
    };
    if (pin) place(pin, {geocode: false});
    map.addListener('click', e => place(toLit(e.latLng)));

    if (searchEl && places.PlaceAutocompleteElement) {
      const ac = new places.PlaceAutocompleteElement({includedRegionCodes: ['kh'], locationBias: SERVICE_BOUNDS});
      ac.className = 'place-search';
      searchEl.replaceChildren(ac);
      const pick = async p => {
        await p.fetchFields({fields: ['id', 'displayName', 'formattedAddress', 'location']});
        if (!p.location) return;
        const pos = round(toLit(p.location));
        map.setCenter(pos); map.setZoom(18);
        place(pos, {geocode: false});
        onPlace?.({...pos, address: p.formattedAddress || '', name: p.displayName || '', placeId: p.id || ''});
      };
      ac.addEventListener('gmp-select', e => e.placePrediction && pick(e.placePrediction.toPlace()));
      ac.addEventListener('gmp-placeselect', e => e.place && pick(e.place)); // versi API lama
    }
    return {
      setPin(pos) { map.setCenter(pos); map.setZoom(18); place(pos); },
    };
  },

  _leafletPicker({mapEl, searchEl, pin, onPin}) {
    if (searchEl) searchEl.replaceChildren();
    if (!window.L) return {setPin() {}};
    const map = L.map(mapEl).setView(pin ? [pin.lat, pin.lng] : [SIHANOUKVILLE.lat, SIHANOUKVILLE.lng], pin ? 17 : 13);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {maxZoom: 19, attribution: '© OpenStreetMap'}).addTo(map);
    let marker = pin ? L.marker([pin.lat, pin.lng]).addTo(map) : null;
    const place = p => {
      p = {lat: Number(p.lat.toFixed(6)), lng: Number(p.lng.toFixed(6))};
      if (marker) marker.setLatLng([p.lat, p.lng]); else marker = L.marker([p.lat, p.lng]).addTo(map);
      onPin(p);
    };
    map.on('click', e => place({lat: e.latlng.lat, lng: e.latlng.lng}));
    return {setPin(p) { map.setView([p.lat, p.lng], 17); place(p); }};
  },

  // markers: tampilkan beberapa titik (admin). points: [{lat,lng,label}]
  async markers(mapEl, points) {
    let ok = false;
    try { ok = await this.load(); } catch (e) { console.warn(e); }
    if (ok) {
      const [{Map, InfoWindow}, {AdvancedMarkerElement}, {LatLngBounds}] = await Promise.all([
        google.maps.importLibrary('maps'), google.maps.importLibrary('marker'), google.maps.importLibrary('core')]);
      const map = new Map(mapEl, {center: points[0], zoom: 16, mapId: this.cfg.googleMapId || 'DEMO_MAP_ID', streetViewControl: false});
      const info = new InfoWindow(), bounds = new LatLngBounds();
      for (const p of points) {
        const m = new AdvancedMarkerElement({map, position: {lat: p.lat, lng: p.lng}, title: p.label});
        m.addListener('click', () => { info.setContent(esc(p.label)); info.open({map, anchor: m}); });
        bounds.extend({lat: p.lat, lng: p.lng});
      }
      if (points.length > 1) map.fitBounds(bounds, 40);
      return;
    }
    if (!window.L) return;
    const map = L.map(mapEl);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {maxZoom: 19, attribution: '© OpenStreetMap'}).addTo(map);
    const pts = points.map(p => (L.marker([p.lat, p.lng]).addTo(map).bindPopup(esc(p.label)), [p.lat, p.lng]));
    map.fitBounds(pts, {maxZoom: 17, padding: [30, 30]});
  },
};
