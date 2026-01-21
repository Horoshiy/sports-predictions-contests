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
  Typography,
  Alert,
} from '@mui/material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import type { ChallengeFormData } from '../../types/challenge.types'

// Validation schema for challenge form
const challengeSchema = z.object({
  opponentId: z.number().min(1, 'Please select an opponent'),
  eventId: z.number().min(1, 'Please select an event'),
  message: z.string().max(500, 'Message cannot exceed 500 characters').optional(),
})

interface ChallengeDialogProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: ChallengeFormData) => void
  availableOpponents?: Array<{ id: number; name: string }>
  availableEvents?: Array<{ id: number; title: string; startDate: string }>
  loading?: boolean
  error?: string | null
}

export const ChallengeDialog: React.FC<ChallengeDialogProps> = ({
  open,
  onClose,
  onSubmit,
  availableOpponents = [],
  availableEvents = [],
  loading = false,
  error = null,
}) => {
  const {
    control,
    handleSubmit,
    reset,
    formState: { errors, isValid },
  } = useForm<ChallengeFormData>({
    resolver: zodResolver(challengeSchema),
    defaultValues: {
      opponentId: 0,
      eventId: 0,
      message: '',
    },
  })

  const handleClose = () => {
    reset()
    onClose()
  }

  const handleFormSubmit = (data: ChallengeFormData) => {
    onSubmit(data)
  }

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>Create Challenge</DialogTitle>
      <form onSubmit={handleSubmit(handleFormSubmit)}>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3, pt: 1 }}>
            {error && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {error}
              </Alert>
            )}

            <Typography variant="body2" color="text.secondary">
              Challenge another user to a head-to-head prediction competition on a specific event.
            </Typography>

            <FormControl fullWidth error={!!errors.opponentId}>
              <InputLabel>Select Opponent</InputLabel>
              <Controller
                name="opponentId"
                control={control}
                render={({ field }) => (
                  <Select
                    {...field}
                    label="Select Opponent"
                    value={field.value || ''}
                    onChange={(e) => field.onChange(Number(e.target.value))}
                  >
                    <MenuItem value="">
                      <em>Choose an opponent</em>
                    </MenuItem>
                    {availableOpponents.map((opponent) => (
                      <MenuItem key={opponent.id} value={opponent.id}>
                        {opponent.name}
                      </MenuItem>
                    ))}
                  </Select>
                )}
              />
              {errors.opponentId && (
                <Typography variant="caption" color="error" sx={{ mt: 0.5 }}>
                  {errors.opponentId.message}
                </Typography>
              )}
            </FormControl>

            <FormControl fullWidth error={!!errors.eventId}>
              <InputLabel>Select Event</InputLabel>
              <Controller
                name="eventId"
                control={control}
                render={({ field }) => (
                  <Select
                    {...field}
                    label="Select Event"
                    value={field.value || ''}
                    onChange={(e) => field.onChange(Number(e.target.value))}
                  >
                    <MenuItem value="">
                      <em>Choose an event</em>
                    </MenuItem>
                    {availableEvents.map((event) => (
                      <MenuItem key={event.id} value={event.id}>
                        <Box>
                          <Typography variant="body2">{event.title}</Typography>
                          <Typography variant="caption" color="text.secondary">
                            {new Date(event.startDate).toLocaleDateString()}
                          </Typography>
                        </Box>
                      </MenuItem>
                    ))}
                  </Select>
                )}
              />
              {errors.eventId && (
                <Typography variant="caption" color="error" sx={{ mt: 0.5 }}>
                  {errors.eventId.message}
                </Typography>
              )}
            </FormControl>

            <Controller
              name="message"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  label="Challenge Message (Optional)"
                  multiline
                  rows={3}
                  placeholder="Add a personal message to your challenge..."
                  error={!!errors.message}
                  helperText={errors.message?.message}
                  fullWidth
                />
              )}
            />

            <Box sx={{ bgcolor: 'grey.50', p: 2, borderRadius: 1 }}>
              <Typography variant="body2" color="text.secondary">
                <strong>Challenge Rules:</strong>
                <br />
                • Your opponent has 24 hours to accept or decline
                • Once accepted, you'll compete on predictions for the selected event
                • Winner is determined by prediction accuracy and points scored
              </Typography>
            </Box>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} disabled={loading}>
            Cancel
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={!isValid || loading}
          >
            {loading ? 'Creating...' : 'Send Challenge'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}

export default ChallengeDialog
