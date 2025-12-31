import { clientAPI } from './index'

export const getClients = (cluster) => clientAPI.getAll(cluster)
export const getClient = (cluster, identity) => clientAPI.getByClusterAndIdentity(cluster, identity)
export const createClient = (data) => clientAPI.create(data)
export const updateClient = (cluster, identity, data) => clientAPI.update(cluster, identity, data)
export const deleteClient = (cluster, identity) => clientAPI.delete(cluster, identity)

