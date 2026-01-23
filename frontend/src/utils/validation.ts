import { z } from 'zod';

export const validateEmail = (email: string): boolean => {
  // More robust email validation
  const re = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return re.test(email);
};

export const validatePassword = (password: string): boolean => {
  // Require at least 8 characters with uppercase, lowercase, and number
  if (password.length < 8) return false;
  const hasUpper = /[A-Z]/.test(password);
  const hasLower = /[a-z]/.test(password);
  const hasNumber = /\d/.test(password);
  return hasUpper && hasLower && hasNumber;
};

export const getPasswordStrength = (password: string): 'weak' | 'medium' | 'strong' => {
  if (password.length < 8) return 'weak';
  
  let strength = 0;
  if (/[a-z]/.test(password)) strength++;
  if (/[A-Z]/.test(password)) strength++;
  if (/\d/.test(password)) strength++;
  if (/[^a-zA-Z0-9]/.test(password)) strength++;
  if (password.length >= 12) strength++;
  
  if (strength <= 2) return 'weak';
  if (strength <= 3) return 'medium';
  return 'strong';
};

// Proper Zod validation schema for contests
export const contestSchema = z.object({
  title: z.string()
    .trim()
    .min(3, 'Title must be at least 3 characters')
    .max(200, 'Title cannot exceed 200 characters'),
  description: z.string()
    .max(1000, 'Description cannot exceed 1000 characters')
    .optional(),
  sportType: z.string()
    .trim()
    .min(1, 'Sport type is required'),
  rules: z.string().optional(),
  startDate: z.date().refine(date => date > new Date(), {
    message: 'Start date must be in the future'
  }),
  endDate: z.date(),
  maxParticipants: z.number()
    .int()
    .min(0, 'Must be 0 or positive')
    .max(10000, 'Cannot exceed 10,000 participants'),
}).refine(data => data.endDate > data.startDate, {
  message: 'End date must be after start date',
  path: ['endDate'],
});

export type ContestFormData = z.infer<typeof contestSchema>;
