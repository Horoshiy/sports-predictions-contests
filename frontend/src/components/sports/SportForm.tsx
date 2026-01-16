import React from 'react'
import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button, Box, Grid } from '@mui/material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { sportSchema, type SportFormData, generateSlug } from '../../utils/sports-validation'
import type { Sport } from '../../types/sports.types'

interface SportFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: SportFormData) => void
  sport?: Sport | null
  loading?: boolean
}

export const SportForm: React.FC<SportFormProps> = ({ open, onClose, onSubmit, sport, loading = false }) => {
  const isEditing = !!sport

  const defaultValues = React.useMemo(() => ({
    name: sport?.name || '',
    slug: sport?.slug || '',
    description: sport?.description || '',
    iconUrl: sport?.iconUrl || '',
  }), [sport])

  const { control, handleSubmit, reset, watch, setValue, formState: { errors, isValid } } = useForm<SportFormData>({
    resolver: zodResolver(sportSchema),
    defaultValues,
    mode: 'onChange',
  })

  const name = watch('name')
  const slugTouched = React.useRef(false)

  React.useEffect(() => {
    if (!isEditing && name && !slugTouched.current) {
      setValue('slug', generateSlug(name))
    }
  }, [name, isEditing, setValue])

  React.useEffect(() => {
    reset(defaultValues)
  }, [defaultValues, reset])

  const handleClose = () => {
    slugTouched.current = false
    reset()
    onClose()
  }

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>{isEditing ? 'Edit Sport' : 'Create Sport'}</DialogTitle>
      <form onSubmit={handleSubmit(onSubmit)}>
        <DialogContent>
          <Box sx={{ mt: 1 }}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Controller
                  name="name"
                  control={control}
                  render={({ field }) => (
                    <TextField {...field} label="Name" fullWidth required error={!!errors.name} helperText={errors.name?.message} disabled={loading} />
                  )}
                />
              </Grid>
              <Grid item xs={12}>
                <Controller
                  name="slug"
                  control={control}
                  render={({ field }) => (
                    <TextField {...field} label="Slug" fullWidth error={!!errors.slug} helperText={errors.slug?.message || 'Auto-generated from name'} disabled={loading} onChange={(e) => { slugTouched.current = true; field.onChange(e) }} />
                  )}
                />
              </Grid>
              <Grid item xs={12}>
                <Controller
                  name="description"
                  control={control}
                  render={({ field }) => (
                    <TextField {...field} label="Description" fullWidth multiline rows={3} error={!!errors.description} helperText={errors.description?.message} disabled={loading} />
                  )}
                />
              </Grid>
              <Grid item xs={12}>
                <Controller
                  name="iconUrl"
                  control={control}
                  render={({ field }) => (
                    <TextField {...field} label="Icon URL" fullWidth error={!!errors.iconUrl} helperText={errors.iconUrl?.message} disabled={loading} />
                  )}
                />
              </Grid>
            </Grid>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} disabled={loading}>Cancel</Button>
          <Button type="submit" variant="contained" disabled={loading || !isValid}>
            {loading ? 'Saving...' : isEditing ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}

export default SportForm
