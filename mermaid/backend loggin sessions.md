```mermaid
flowchart TD
    A[User Login Request] --> B{Validate Google JWT}
    
    B -->|Invalid| C[Return Auth Error]
    B -->|Valid| D[Extract Key JWT Parameters]
    
    D --> E[Extract sub as UserID]
    D --> F[Verify iss is accounts.google.com]
    D --> G[Validate aud matches client ID]
    D --> H[Check exp for token expiration]
    D --> I[Verify email_verified is true]
    D --> J[Optional: Check hd domain]
    
    E & F & G & H & I & J --> K{User exists in DB?}
    
    K -->|No| L[Create New User Record]
    K -->|Yes| M[Update User Last Login]
    
    L --> N[Generate Session]
    M --> N
    
    N --> O[Create Session Record in DB/Redis]
    
    subgraph Session Creation
        O --> P[Store User ID from sub]
        O --> Q[Set Expiry from exp claim]
        O --> R[Record Token Fingerprint]
        O --> S[Store Original nonce]
        O --> T[Add Client IP & User Agent]
    end
    
    subgraph Session Cookie
        P & Q & R & S & T --> U[Generate Session ID]
        U --> V[Set HttpOnly Cookie]
        V --> W[Set Secure flag]
        V --> X[Set SameSite=Lax]
        V --> Y[Set Path=/]
        V --> Z[Set Domain appropriately]
    end
    
    X & Y & Z --> AA[Return Success Response]
    
    AB[Subsequent Request] --> AC{Check Session Cookie}
    AC -->|Missing| AD[Redirect to Login]
    AC -->|Present| AE[Lookup Session in DB/Redis]
    
    AE -->|Not Found| AF[Session Invalid]
    AF --> AD
    
    AE -->|Found| AG{Validate Session}
    AG -->|Expired| AH[Clear Cookie & Redirect to Login]
    AG -->|Valid| AI[Refresh Session Timeout]
    AI --> AJ[Process Authenticated Request]
    ```