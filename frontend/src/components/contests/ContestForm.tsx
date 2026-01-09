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
  InputLabel,
  Select,
  MenuItem,
  Grid,
} from '@mui/material'
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider'
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { contestSchema, type ContestFormData } from '../../utils/validation'
import type { Contest } from '../../types/contest.types'

interface ContestFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: ContestFormData) => void
  contest?: Contest | null
  loading?: boolean
}

const sportTypes = [
  'Football',
  'Basketball',
  'Baseball',
  'Soccer',
  'Tennis',
  'Hockey',
  'Golf',
  'Boxing',
  'MMA',
  'Other',
]

export const ContestForm: React.FC<ContestFormProps> = ({
  open,
  onClose,
  onSubmit,
  contest,
  loading = false,
}) => {
  const isEditing = !!contest

  const defaultValues = React.useMemo(() => ({
    title: contest?.title || '',
    description: contest?.description || '',
    sportType: contest?.sportType || '',
    rules: contest?.rules || '',
    startDate: contest?.startDate ? new Date(contest.startDate) : new Date(),
    endDate: contest?.endDate ? new Date(contest.endDate) : new Date(),
    maxParticipants: contest?.maxParticipants || 0,
  }), [contest])

  const {
    control,
    handleSubmit,
    reset,
    formState: { errors, isValid },
  } = useForm<ContestFormData>({
    resolver: zodResolver(contestSchema),
    defaultValues,
    mode: 'onChange',
  })

  const handleFormSubmit = (data: ContestFormData) => {
    onSubmit(data)
  }

  const handleClose = () => {
    reset()
    onClose()
  }

  React.useEffect(() => {
    reset(defaultValues)
  }, [defaultValues, reset])

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>
        {isEditing ? 'Edit Contest' : 'Create New Contest'}
      </DialogTitle>
      
      <form onSubmit={handleSubmit(handleFormSubmit)}>
        <DialogContent>
          <Box sx={{ mt: 1 }}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Controller
                  name="title"
                  control={control}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      label="Contest Title"
                      fullWidth
                      required
                      error={!!errors.title}
                      helperText={errors.title?.message}
                      disabled={loading}
                    />
                  )}
                />
              </Grid>

              <Grid item xs={12}>
                <Controller
                  name="description"
                  control={control}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      label="Description"
                      fullWidth
                      multiline
                      rows={3}
                      error={!!errors.description}
                      helperText={errors.description?.message}
                      disabled={loading}
                    />
                  )}
                />
              </Grid>

              <Grid item xs={12} sm={6}>
                <Controller
                  name="sportType"
                  control={control}
                  render={({ field }) => (
                    <FormControl fullWidth required error={!!errors.sportType}>
                      <InputLabel>Sport Type</InputLabel>
                      <Select
                        {...field}
                        label="Sport Type"
                        disabled={loading}
                      >
                        {sportTypes.map((sport) => (
                          <MenuItem key={sport} value={sport}>
                            {sport}
                          </MenuItem>
                        ))}
                      </Select>
                      {errors.sportType && (
                        <Box sx={{ color: 'error.main', fontSize: '0.75rem', mt: 0.5 }}>
                          {errors.sportType.message}
                        </Box>
                      )}
                    </FormControl>
                  )}
                />
              </Grid>

              <Grid item xs={12} sm={6}>
                <Controller
                  name="maxParticipants"
                  control={control}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      label="Max Participants (0 = unlimited)"
                      type="number"
                      fullWidth
                      error={!!errors.maxParticipants}
                      helperText={errors.maxParticipants?.message}
                      disabled={loading}
                      onChange={(e) => field.onChange(parseInt(e.target.value) || 0)}
                    />
                  )}
                />
              </Grid>

              <Grid item xs={12} sm={6}>
                <LocalizationProvider dateAdapter={AdapterDateFns}>
                  <Controller
                    name="startDate"
                    control={control}
                    render={({ field }) => (
                      <DateTimePicker
                        {...field}
                        label="Start Date & Time"
                        disabled={loading}
                        slotProps={{
                          textField: {
                            fullWidth: true,
                            error: !!errors.startDate,
                            helperText: errors.startDate?.message,
                          },
                        }}
                      />
                    )}
                  />
                </LocalizationProvider>
              </Grid>

              <Grid item xs={12} sm={6}>
                <LocalizationProvider dateAdapter={AdapterDateFns}>
                  <Controller
                    name="endDate"
                    control={control}
                    render={({ field }) => (
                      <DateTimePicker
                        {...field}
                        label="End Date & Time"
                        disabled={loading}
                        slotProps={{
                          textField: {
                            fullWidth: true,
                            error: !!errors.endDate,
                            helperText: errors.endDate?.message,
                          },
                        }}
                      />
                    )}
                  />
                </LocalizationProvider>
              </Grid>

              <Grid item xs={12}>
                <Controller
                  name="rules"
                  control={control}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      label="Contest Rules (JSON format)"
                      fullWidth
                      multiline
                      rows={4}
                      error={!!errors.rules}
                      helperText={errors.rules?.message || 'Enter contest rules in JSON format'}
                      disabled={loading}
                    />
                  )}
                />
              </Grid>
            </Grid>
          </Box>
        </DialogContent>

        <DialogActions>
          <Button onClick={handleClose} disabled={loading}>
            Cancel
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={loading || !isValid}
          >
            {loading ? 'Saving...' : isEditing ? 'Update Contest' : 'Create Contest'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}

export default ContestForm
