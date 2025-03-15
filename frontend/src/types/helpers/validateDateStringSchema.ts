import { z } from 'zod'

export const dateStringSchema = z.preprocess((val) => {
  if (typeof val === 'string') {
    return new Date(val)
  }
  return val
}, z.date())

export const correlationIDSchema = z.number().and(
  z.object({
    id: z.literal(true),
  })
)
