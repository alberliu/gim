import { createGrpcWebTransport } from '@connectrpc/connect-web'
import { createClient, type Interceptor } from '@connectrpc/connect'
import { auth } from '@/store/auth'
import { UserExtService } from '@/gen/proto/business/user.ext_pb'
import { FriendExtService } from '@/gen/proto/business/friend.ext_pb'
import { GroupExtService } from '@/gen/proto/business/group.ext_pb'
import { MessageExtService as BizMessageExtService } from '@/gen/proto/business/message.ext_pb'
import { MessageExtService as LogicMessageExtService } from '@/gen/proto/logic/message.ext_pb'

export const GRPC_BASE_URL = 'http://127.0.0.1:8080'
export const UPLOAD_URL = 'http://127.0.0.1:8081/upload'

const tokenInterceptor: Interceptor = (next) => async (req) => {
  const t = auth.token
  const uid = auth.userId
  const did = auth.deviceId
  if (t) req.header.set('token', t)
  if (uid) req.header.set('user_id', uid)
  if (did) req.header.set('device_id', did)
  return await next(req)
}

const transport = createGrpcWebTransport({
  baseUrl: GRPC_BASE_URL,
  interceptors: [tokenInterceptor],
})

export const userClient = createClient(UserExtService, transport)
export const friendClient = createClient(FriendExtService, transport)
export const groupClient = createClient(GroupExtService, transport)
export const messageClient = createClient(BizMessageExtService, transport)
export const messageLogicClient = createClient(LogicMessageExtService, transport)

;(window as any).__gim = {
  userClient,
  friendClient,
  groupClient,
  messageClient,
  messageLogicClient,
}
