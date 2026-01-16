import { z } from 'zod'

// Main prediction form schema
export const predictionSchema = z.object({
  eventId: z.number().min(1, 'Event selection is required'),
  predictionType: z.enum(['winner', 'score', 'combined', 'props'], {
    errorMap: () => ({ message: 'Select prediction type' }),
  }),
  winner: z.enum(['home', 'away', 'draw']).optional(),
  homeScore: z.number().min(0, 'Score must be 0 or greater').optional(),
  awayScore: z.number().min(0, 'Score must be 0 or greater').optional(),
}).refine(
  (data) => {
    if (data.predictionType === 'winner' || data.predictionType === 'combined') {
      return data.winner !== undefined
    }
    return true
  },
  { message: 'Winner selection is required', path: ['winner'] }
).refine(
  (data) => {
    if (data.predictionType === 'score' || data.predictionType === 'combined') {
      return data.homeScore !== undefined && data.awayScore !== undefined
    }
    return true
  },
  { message: 'Scores are required', path: ['homeScore'] }
)

export type PredictionFormData = z.infer<typeof predictionSchema>

// Single prop prediction schema
export const propPredictionSchema = z.object({
  propTypeId: z.number().min(1, 'Prop type is required'),
  propSlug: z.string().min(1),
  line: z.number().optional(),
  selection: z.string().min(1, 'Selection is required'),
  playerId: z.string().optional(),
  pointsValue: z.number().default(2),
})

export type PropPredictionFormData = z.infer<typeof propPredictionSchema>

// Helper to convert form data to JSON string
export const formDataToPredictionData = (data: PredictionFormData): string => {
  const result: Record<string, unknown> = {}
  
  if (data.predictionType === 'winner' || data.predictionType === 'combined') {
    result.winner = data.winner
  }
  
  if (data.predictionType === 'score' || data.predictionType === 'combined') {
    result.homeScore = data.homeScore
    result.awayScore = data.awayScore
  }
  
  return JSON.stringify(result)
}

// Helper to convert props form data to JSON
export const propsFormDataToPredictionData = (props: PropPredictionFormData[]): string => {
  return JSON.stringify({
    type: 'props',
    props: props.map(p => ({
      prop_type_id: p.propTypeId,
      prop_slug: p.propSlug,
      line: p.line,
      selection: p.selection,
      player_id: p.playerId,
      points_value: p.pointsValue,
    })),
  })
}

// Helper to parse prediction data JSON to form data
export const predictionDataToFormData = (
  predictionData: string,
  eventId: number
): Partial<PredictionFormData> => {
  try {
    const parsed = JSON.parse(predictionData)
    const hasWinner = 'winner' in parsed
    const hasScore = 'homeScore' in parsed && 'awayScore' in parsed
    
    let predictionType: 'winner' | 'score' | 'combined' | 'props' = 'winner'
    if (parsed.type === 'props') predictionType = 'props'
    else if (hasWinner && hasScore) predictionType = 'combined'
    else if (hasScore) predictionType = 'score'
    
    return {
      eventId,
      predictionType,
      winner: parsed.winner,
      homeScore: parsed.homeScore,
      awayScore: parsed.awayScore,
    }
  } catch {
    return { eventId, predictionType: 'winner' }
  }
}
