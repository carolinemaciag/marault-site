# Chatbot Testing Guide - Alternative Questions & Context Memory

## Quick Test Cases

### 1. Alternative Question Wording Tests

#### Data Visibility Audit
```
User: "audit my data"
Expected: Full explanation of what audit is
Bot will recognize: New keyword variation

User: "understand data visibility"
Expected: Clear explanation + company context
Bot will recognize: New keyword variation

User: "how does audit work?"
Expected: Process explanation
Bot will recognize: New keyword variation
```

#### Executive Dashboard
```
User: "tell me about dashboards"
Expected: Dashboard explanation with Marault positioning
Bot will recognize: New keyword variation

User: "dashboard meaning"
Expected: Clear definition + why it matters
Bot will recognize: New keyword variation
```

#### Getting Started
```
User: "let's talk"
Expected: Starting process with timeline
Bot will recognize: New keyword "let's talk"

User: "move forward"
Expected: Action steps with discovery call detail
Bot will recognize: New keyword variation

User: "want help"
Expected: Getting started process
Bot will recognize: New keyword variation
```

---

### 2. Context Memory Tests

#### Test: Audit → Tell Me More
```
Message 1 (User): "What is a data visibility audit?"
Bot Response: [Full audit explanation]

Message 2 (User): "tell me more"
Bot Response: [ELABORATED] Details 6 deliverables:
- KPI definitions
- Clean reporting baseline
- Data inventory and lineage
- Data quality scorecard
- Executive dashboard
- Stabilization roadmap
Includes: "At Marault Intelligence, we combine strategy and technology..."
```

#### Test: Dashboard → How Do We Start
```
Message 1 (User): "What is an executive dashboard?"
Bot Response: [Dashboard explanation]

Message 2 (User): "How do we start?"
Bot Response: [SERVICE-SPECIFIC] 
- KPI alignment conversation
- Design, build, validate phases
- Team training included
Includes: "At Marault Intelligence, we ensure the dashboard becomes a real tool..."
```

#### Test: Website → What Else
```
Message 1 (User): "Tell me about custom website builds"
Bot Response: [Website explanation]

Message 2 (User): "What else should we think about?"
Bot Response: [UPSELL] Suggests:
- If web focus: Data analytics services
- If both needed: Strategic partnership approach
Includes: "At Marault Intelligence, we help companies compete through clarity..."
```

---

### 3. Company/Brand Emphasis Tests

#### Test: Brand References
Look for these phrases in responses:
- "At Marault Intelligence, we..." ✓
- "We specialize in..." ✓
- "We've worked across industries..." ✓
- Specific metrics: "50% reduction" ✓
- Industry mentions: "Finance, SaaS, E-Commerce..." ✓

#### Test: Clarity vs Generic
```
OLD: "Visit inquire page"
NEW: "Visit inquire page—our team responds within 24 hours"

OLD: "We offer 8 core services"
NEW: "Marault Intelligence specializes in transforming data and design into competitive advantages"

OLD: "We design dashboards"
NEW: "At Marault Intelligence, we design dashboards that actually get used and drive action"
```

---

### 4. New Context-Aware Handler Tests

#### Test: "What Else" Handler
```
User has asked about: Data Audit
User now asks: "What else?"
Bot should: Recommend Dashboard + Forecasting + Revenue Analytics
Include: "At Marault Intelligence, we help companies build incrementally"
✓ Recognize context
✓ Suggest related services
✓ Emphasize strategic progression
```

#### Test: Industry/Company Handler
```
User: "We're a SaaS company, is this relevant?"
Bot should:
✓ Recognize industry-specific question
✓ Mention SaaS experience
✓ Emphasize customization
✓ Invite deeper conversation
```

#### Test: Comparison Handler
```
User: "Why choose Marault over other agencies?"
Bot should:
✓ Emphasize strategic + technical approach
✓ Mention "not just building, enabling"
✓ Focus on measurable outcomes
✓ Invite detailed conversation
```

#### Test: Timeline Handler
```
User: "How fast can you do this?"
Bot should:
✓ Provide specific timelines (2-3 weeks audit, 4-8 weeks dashboard, etc.)
✓ Emphasize quality over speed
✓ Note efficiency
✓ Invite discussion of specific timeline
```

---

### 5. Follow-Up Elaboration Tests

#### Scenario: Results/ROI Question
```
Context: User asked about audit
User now: "What results can we expect?"
Bot should:
✓ Recognize context = audit
✓ Provide audit-specific outcomes:
  - Clearer decisions in 2-3 weeks
  - 50%+ faster review cycles
  - No conflicting reports
  - Better confidence
  - Ongoing roadmap
✓ Include: "At Marault Intelligence, we measure success by..."
```

#### Scenario: Implementation Details
```
Context: User asked about dashboard
User now: "What exactly happens?"
Bot should:
✓ Recognize context = dashboard
✓ Provide step-by-step process
✓ Mention KPI alignment conversation
✓ Reference training + governance
✓ Include metrics improvement mention (50%+ time reduction)
```

---

### 6. Edge Case Tests

#### Test: Multiple Services in History
```
Context: "I asked about audit, then dashboard, now website"
User: "What else?"
Bot should:
✓ Recognize most recent context (website)
✓ Suggest related web services OR
✓ Suggest analytics for complete strategy
✓ Mention strategic partnership
```

#### Test: No Context
```
User (first message): "Tell me more"
Bot should:
✓ Recognize no previous context
✓ Default to general "tell me more" → services overview
✓ Invite them to ask about specific service
```

#### Test: Unclear Context
```
Context: Generic "what do you do" discussion
User: "Tell me more"
Bot should:
✓ Not assume specific service
✓ Default to company mission statement
✓ Invite specific service question
```

---

## Real Conversation Examples

### Example 1: Complete User Journey
```
User: "audit my data"
Bot: ✓ Recognizes new keyword, explains audit

User: "tell me more"
Bot: ✓ Remembers context, elaborates with 6 deliverables

User: "What happens after?"
Bot: ✓ Context-aware, explains implementation roadmap

User: "How do we start?"
Bot: ✓ Service-specific, explains discovery call process

User: "What else?"
Bot: ✓ Recommends Dashboard + Forecasting
```

### Example 2: Cross-Service Discovery
```
User: "We need better reporting"
Bot: ✓ Recognizes need, explains dashboards

User: "How is this different?"
Bot: ✓ Compares to operational dashboards

User: "What about our data?"
Bot: ✓ Recommends audit first, then dashboard

User: "Can you do both?"
Bot: ✓ Explains progression strategy
```

### Example 3: Industry-Specific
```
User: "We're a SaaS company"
Bot: ✓ Recognizes industry context

User: "What services are relevant?"
Bot: ✓ Emphasizes SaaS-specific metrics

User: "Tell me more"
Bot: ✓ Elaborates with SaaS examples
```

---

## Verification Checklist

- [ ] Alternative keywords recognize new question phrasings
- [ ] Context memory works for "tell me more" requests
- [ ] Context memory works for "how do we start" requests
- [ ] Context memory works for "what else" requests
- [ ] Service-specific elaborations appear based on history
- [ ] Marault Intelligence branding appears in 40+ responses
- [ ] Industry/company context recognized and respected
- [ ] ROI/results clearly explained per service
- [ ] Timeline information provided when asked
- [ ] All responses end with clear next step
- [ ] No emojis in any response
- [ ] No forward slashes (uses page references instead)
- [ ] Company expertise emphasized throughout
- [ ] New context-aware handlers functioning
- [ ] Default response helpful and brand-appropriate

---

## Common Testing Commands

```
# Test alternative wording
"audit my data"
"tell me about dashboards"  
"let's talk"
"move forward"
"want help"

# Test context memory
First: "What is a data audit?"
Then: "tell me more"

# Test company emphasis
Look for: "At Marault Intelligence"
Look for: "we specialize"
Look for: "50% reduction"

# Test new handlers
"We're in finance"
"Why Marault?"
"How fast?"
"What else?"

# Test elaboration
"What results?"
"What happens?"
"Give me details"
"Expand on that"
```

---

## Success Metrics

✓ Users find answers even with alternative wording
✓ Follow-up questions get deeper, contextualized responses
✓ Marault's expertise is evident in every response
✓ Users understand the service progression (audit → dashboard → forecasting)
✓ Users feel informed and confident moving to inquire page
✓ No generic or repetitive responses
✓ Brand voice consistent throughout
✓ Response quality improves with conversation history

