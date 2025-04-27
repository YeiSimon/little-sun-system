// google-identity.d.ts

// Define the response returned by Google after a sign-in
interface GoogleCredentialResponse {
    credential: string;
    select_by: string;
  }
  
  // Options for initializing Google Identity Services
  interface GoogleIdentityServiceOptions {
    client_id: string;
    callback: (response: GoogleCredentialResponse) => void;
    auto_select?: boolean;
    cancel_on_tap_outside?: boolean;
    context?: 'signin' | 'signup';
  }
  
  // Options for rendering the Google Sign-In button
  interface RenderButtonOptions {
    type?: 'standard' | 'icon';
    theme?: 'outline' | 'filled_blue' | 'filled_black';
    size?: 'small' | 'medium' | 'large';
    text?: 'signin_with' | 'signup_with' | 'continue_with';
    shape?: 'rectangular' | 'pill' | 'circle' | 'square';
    logo_alignment?: 'left' | 'center';
  }
  
  // Define the API available under google.accounts.id
  interface GoogleId {
    initialize(options: GoogleIdentityServiceOptions): void;
    renderButton(container: HTMLElement, options?: RenderButtonOptions): void;
    prompt(callback?: (notification: any) => void): void;
  }
  
  // Define the GoogleAccounts interface which contains the id property
  interface GoogleAccounts {
    id: GoogleId;
  }
  
  // Declare the global 'google' variable to match the API structure.
  declare const google: {
    accounts: GoogleAccounts;
  };
  