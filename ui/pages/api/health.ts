import type { NextApiRequest, NextApiResponse } from 'next'

type Health = { status: 'ok' }

export default function handler(
  _req: NextApiRequest,
  res: NextApiResponse<Health>
) {
  res.status(200).json({ status: 'ok' })
}

