import React from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box,
  FormControl,
  FormLabel,
  RadioGroup,
  FormControlLabel,
  Radio,
  Typography,
  Divider,
  Grid,
  FormHelperText,
  Alert,
} from '@mui/material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { 
  predictionSchema, 
  type PredictionFormData,
  type PropPredictionFormData,
  predictionDataToFormData,
  formDataToPredictionData,
  propsFormDataToPredictionData,
} from '../../utils/prediction-validation'
import type { Event, Prediction } from '../../types/prediction.types'
import { formatDate } from '../../utils/date-utils'
import { usePropTypes, usePotentialCoefficient } from '../../hooks/use-predictions'
import { PropTypeSelector } from './PropTypeSelector'
import { CoefficientIndicator } from './CoefficientIndicator'

interface PredictionFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (predictionData: string) => void
  event: Event | null
  prediction?: Prediction | null
  loading?: boolean
}

export const PredictionForm: React.FC<PredictionFormProps> = ({
  open,
  onClose,
  onSubmit,
  event,
  prediction,
  loading = false,
}) => {
  const isEditing = !!prediction
  const [selectedProps, setSelectedProps] = React.useState<PropPredictionFormData[]>([])
  const [propsError, setPropsError] = React.useState<string | null>(null)
  const { data: propTypes = [] } = usePropTypes(event?.sportType || '')
  const { data: coefficientData } = usePotentialCoefficient(event?.id)

  const defaultValues = React.useMemo(() => {
    if (prediction && event) {
      return predictionDataToFormData(prediction.predictionData, event.id)
    }
    return {
      eventId: event?.id || 0,
      predictionType: 'winner' as const,
      winner: undefined,
      homeScore: undefined,
      awayScore: undefined,
    }
  }, [prediction, event])

  const {
    control,
    handleSubmit,
    watch,
    reset,
    formState: { errors },
  } = useForm<PredictionFormData>({
    resolver: zodResolver(predictionSchema),
    defaultValues,
    mode: 'onBlur',
  })

  const predictionType = watch('predictionType')

  React.useEffect(() => {
    if (open) {
      reset(defaultValues)
      setSelectedProps([])
      setPropsError(null)
    }
  }, [open, defaultValues, reset])

  const handleFormSubmit = (data: PredictionFormData) => {
    let predictionData: string
    if (data.predictionType === 'props') {
      if (selectedProps.length === 0) {
        setPropsError('Please add at least one prop prediction')
        return
      }
      const invalidProps = selectedProps.filter(p => !p.selection)
      if (invalidProps.length > 0) {
        setPropsError('Please make a selection for all props')
        return
      }
      setPropsError(null)
      predictionData = propsFormDataToPredictionData(selectedProps)
    } else {
      predictionData = formDataToPredictionData(data)
    }
    onSubmit(predictionData)
  }

  const handleClose = () => {
    reset()
    setSelectedProps([])
    setPropsError(null)
    onClose()
  }

  if (!event) return null

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        {isEditing ? 'Edit Prediction' : 'Make Prediction'}
      </DialogTitle>
      
      <form onSubmit={handleSubmit(handleFormSubmit)}>
        <DialogContent>
          <Box sx={{ mb: 3, p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
            <Typography variant="subtitle2" color="text.secondary">
              {event.sportType}
            </Typography>
            <Typography variant="h6">{event.title}</Typography>
            <Box sx={{ mt: 1, textAlign: 'center' }}>
              <Typography variant="body1" fontWeight="medium">
                {event.homeTeam} vs {event.awayTeam}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {formatDate(event.eventDate)}
              </Typography>
            </Box>
          </Box>

          <Divider sx={{ mb: 3 }} />

          {coefficientData && coefficientData.coefficient > 1 && (
            <Box sx={{ mb: 2 }}>
              <CoefficientIndicator
                coefficient={coefficientData.coefficient}
                tier={coefficientData.tier}
                hoursUntilEvent={coefficientData.hoursUntilEvent}
              />
            </Box>
          )}

          <Controller
            name="predictionType"
            control={control}
            render={({ field }) => (
              <FormControl component="fieldset" sx={{ mb: 3 }}>
                <FormLabel>Prediction Type</FormLabel>
                <RadioGroup {...field} row>
                  <FormControlLabel value="winner" control={<Radio />} label="Winner Only" />
                  <FormControlLabel value="score" control={<Radio />} label="Exact Score" />
                  <FormControlLabel value="combined" control={<Radio />} label="Both" />
                  {propTypes.length > 0 && (
                    <FormControlLabel value="props" control={<Radio />} label="Props" />
                  )}
                </RadioGroup>
              </FormControl>
            )}
          />

          {predictionType === 'props' && propTypes.length > 0 && (
            <>
              {propsError && (
                <Alert severity="error" sx={{ mb: 2 }}>{propsError}</Alert>
              )}
              <PropTypeSelector
                propTypes={propTypes}
                selectedProps={selectedProps}
                onPropsChange={(props) => {
                  setSelectedProps(props)
                  setPropsError(null)
                }}
                homeTeam={event.homeTeam}
                awayTeam={event.awayTeam}
                disabled={loading}
              />
            </>
          )}

          {(predictionType === 'winner' || predictionType === 'combined') && (
            <Controller
              name="winner"
              control={control}
              render={({ field }) => (
                <FormControl component="fieldset" sx={{ mb: 3 }} error={!!errors.winner}>
                  <FormLabel>Select Winner</FormLabel>
                  <RadioGroup {...field} value={field.value || ''}>
                    <FormControlLabel value="home" control={<Radio />} label={event.homeTeam} />
                    <FormControlLabel value="away" control={<Radio />} label={event.awayTeam} />
                    <FormControlLabel value="draw" control={<Radio />} label="Draw" />
                  </RadioGroup>
                  {errors.winner && (
                    <FormHelperText>{errors.winner.message}</FormHelperText>
                  )}
                </FormControl>
              )}
            />
          )}

          {(predictionType === 'score' || predictionType === 'combined') && (
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Controller
                  name="homeScore"
                  control={control}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      label={`${event.homeTeam} Score`}
                      type="number"
                      fullWidth
                      error={!!errors.homeScore}
                      helperText={errors.homeScore?.message}
                      disabled={loading}
                      onChange={(e) => {
                        const val = e.target.value
                        field.onChange(val === '' ? undefined : parseInt(val))
                      }}
                      value={field.value ?? ''}
                      inputProps={{ min: 0 }}
                    />
                  )}
                />
              </Grid>
              <Grid item xs={6}>
                <Controller
                  name="awayScore"
                  control={control}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      label={`${event.awayTeam} Score`}
                      type="number"
                      fullWidth
                      error={!!errors.awayScore}
                      helperText={errors.awayScore?.message}
                      disabled={loading}
                      onChange={(e) => {
                        const val = e.target.value
                        field.onChange(val === '' ? undefined : parseInt(val))
                      }}
                      value={field.value ?? ''}
                      inputProps={{ min: 0 }}
                    />
                  )}
                />
              </Grid>
            </Grid>
          )}
        </DialogContent>

        <DialogActions>
          <Button onClick={handleClose} disabled={loading}>
            Cancel
          </Button>
          <Button type="submit" variant="contained" disabled={loading}>
            {loading ? 'Saving...' : isEditing ? 'Update Prediction' : 'Submit Prediction'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}

export default PredictionForm
