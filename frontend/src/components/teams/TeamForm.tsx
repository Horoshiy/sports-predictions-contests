import React, { useEffect } from 'react'
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Box } from '@mui/material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { teamSchema, type TeamSchemaType } from '../../utils/team-validation'
import { useCreateTeam, useUpdateTeam } from '../../hooks/use-teams'
import type { Team } from '../../types/team.types'

interface TeamFormProps {
  open: boolean
  onClose: () => void
  team?: Team | null
}

export const TeamForm: React.FC<TeamFormProps> = ({ open, onClose, team }) => {
  const isEdit = !!team
  const createTeamMutation = useCreateTeam()
  const updateTeamMutation = useUpdateTeam()

  const { control, handleSubmit, reset, formState: { errors } } = useForm<TeamSchemaType>({
    resolver: zodResolver(teamSchema),
    defaultValues: { name: '', description: '', maxMembers: 10 },
  })

  useEffect(() => {
    if (team) {
      reset({ name: team.name, description: team.description, maxMembers: team.maxMembers })
    } else {
      reset({ name: '', description: '', maxMembers: 10 })
    }
  }, [team, reset])

  const onSubmit = async (data: TeamSchemaType) => {
    try {
      if (isEdit && team) {
        await updateTeamMutation.mutateAsync({ id: team.id, ...data })
      } else {
        await createTeamMutation.mutateAsync(data)
      }
      onClose()
    } catch (error) {
      // Error handled by mutation
    }
  }

  const isPending = createTeamMutation.isPending || updateTeamMutation.isPending

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{isEdit ? 'Edit Team' : 'Create Team'}</DialogTitle>
      <form onSubmit={handleSubmit(onSubmit)}>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            <Controller
              name="name"
              control={control}
              render={({ field }) => (
                <TextField {...field} label="Team Name" required fullWidth error={!!errors.name} helperText={errors.name?.message} />
              )}
            />
            <Controller
              name="description"
              control={control}
              render={({ field }) => (
                <TextField {...field} label="Description" multiline rows={3} fullWidth error={!!errors.description} helperText={errors.description?.message} />
              )}
            />
            <Controller
              name="maxMembers"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  type="number"
                  label="Max Members"
                  fullWidth
                  error={!!errors.maxMembers}
                  helperText={errors.maxMembers?.message || 'Between 2 and 50'}
                  onChange={(e) => field.onChange(parseInt(e.target.value) || 10)}
                />
              )}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Cancel</Button>
          <Button type="submit" variant="contained" disabled={isPending}>
            {isPending ? 'Saving...' : isEdit ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}

export default TeamForm
