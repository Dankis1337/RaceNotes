const DB_NAME = 'racenotes-offline'
const DB_VERSION = 1

let dbInstance = null

function openDb() {
  if (dbInstance) return Promise.resolve(dbInstance)
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION)
    request.onupgradeneeded = (e) => {
      const db = e.target.result
      if (!db.objectStoreNames.contains('races')) {
        db.createObjectStore('races', { keyPath: 'id' })
      }
      if (!db.objectStoreNames.contains('setups')) {
        db.createObjectStore('setups', { keyPath: 'id' })
      }
      if (!db.objectStoreNames.contains('pendingSync')) {
        const store = db.createObjectStore('pendingSync', { keyPath: 'localId', autoIncrement: true })
        store.createIndex('type', 'type')
      }
    }
    request.onsuccess = () => {
      dbInstance = request.result
      dbInstance.onclose = () => { dbInstance = null }
      resolve(dbInstance)
    }
    request.onerror = () => reject(request.error)
  })
}

async function getAll(storeName) {
  const db = await openDb()
  return new Promise((resolve, reject) => {
    const tx = db.transaction(storeName, 'readonly')
    const store = tx.objectStore(storeName)
    const req = store.getAll()
    req.onsuccess = () => resolve(req.result)
    req.onerror = () => reject(req.error)
  })
}

async function putAll(storeName, items) {
  const db = await openDb()
  const tx = db.transaction(storeName, 'readwrite')
  const store = tx.objectStore(storeName)
  store.clear()
  items.forEach(item => store.put(item))
  return new Promise((resolve, reject) => {
    tx.oncomplete = () => resolve()
    tx.onerror = () => reject(tx.error)
  })
}

async function addPendingSync(action) {
  const db = await openDb()
  const tx = db.transaction('pendingSync', 'readwrite')
  tx.objectStore('pendingSync').add({ ...action, createdAt: Date.now() })
  return new Promise((resolve, reject) => {
    tx.oncomplete = () => resolve()
    tx.onerror = () => reject(tx.error)
  })
}

async function getPendingSync() {
  return getAll('pendingSync')
}

async function clearPendingSync() {
  const db = await openDb()
  const tx = db.transaction('pendingSync', 'readwrite')
  tx.objectStore('pendingSync').clear()
  return new Promise((resolve, reject) => {
    tx.oncomplete = () => resolve()
    tx.onerror = () => reject(tx.error)
  })
}

async function removePendingItem(localId) {
  const db = await openDb()
  const tx = db.transaction('pendingSync', 'readwrite')
  tx.objectStore('pendingSync').delete(localId)
  return new Promise((resolve, reject) => {
    tx.oncomplete = () => resolve()
    tx.onerror = () => reject(tx.error)
  })
}

export const offlineDb = {
  cacheRaces: (races) => putAll('races', races),
  cacheSetups: (setups) => putAll('setups', setups),
  getCachedRaces: () => getAll('races'),
  getCachedSetups: () => getAll('setups'),
  addPendingSync,
  getPendingSync,
  clearPendingSync,
  removePendingItem,
}
