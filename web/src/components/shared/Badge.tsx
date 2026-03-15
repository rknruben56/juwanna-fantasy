import './Badge.css';

interface BadgeProps {
  text: string;
  variant?: 'gold' | 'teal' | 'red' | 'muted' | 'default';
  size?: 'sm' | 'md';
}

export function Badge({ text, variant = 'default', size = 'md' }: BadgeProps) {
  return (
    <span className={`badge badge--${variant} badge--${size}`}>
      {text}
    </span>
  );
}
