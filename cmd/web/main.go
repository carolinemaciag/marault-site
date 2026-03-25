package main

import (
	"os"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"strings"
)

func isIPad(r *http.Request) bool {
	ua := strings.ToLower(r.UserAgent())
	return strings.Contains(ua, "ipad")
}

func isMobile(r *http.Request) bool {
	ua := strings.ToLower(r.UserAgent())
	mobileSignals := []string{
		"iphone",
		"android",
		"mobile",
		"ipod",
	}
	for _, s := range mobileSignals {
		if strings.Contains(ua, s) {
			return true
		}
	}
	return false
}

func getBaseTemplate(r *http.Request) string {
	if r.URL.Query().Get("device") == "ipad" {
		return "./internal/templates/base-ipad.html"
	}
	if isMobile(r) {
		return "./internal/templates/base-mobile.html"
	}
	return "./internal/templates/base.html"
}

/*
	=========================
	  GENERIC SERVICE PAGE HANDLER

=========================
*/
func servicePageHandler(templateFile string, title string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(
			getBaseTemplate(r),
			"./internal/templates/"+templateFile,
			"./internal/templates/chatbot.html",
		)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
			Page  string
		}{
			Title: title,
			Page:  "services",
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println(err)
		}
	}
}

/*
	=========================
	  CHATBOT API

=========================
*/
type ChatMessage struct {
	Message string `json:"message"`
	History []struct {
		Sender  string `json:"sender"`
		Message string `json:"message"`
	} `json:"history"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg ChatMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Generate response based on user message and conversation history
	reply := generateChatResponse(msg.Message, msg.History)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{Reply: reply})
}

func generateChatResponse(userMessage string, history []struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}) string {
	userLower := strings.ToLower(userMessage)

	// Build conversation context for memory
	var conversationContext string
	if len(history) > 0 {
		// Get last 3 exchanges for context (to understand conversation flow)
		for i := len(history) - 1; i >= 0 && i >= len(history)-3; i-- {
			conversationContext += strings.ToLower(history[i].Message) + " "
		}
	}

	// CASUAL GREETING HANDLER
	if contains(userLower, []string{"whats up", "what's up", "sup", "yo", "hey", "hi", "hello", "greetings", "good morning", "good afternoon", "good evening"}) && len(userMessage) < 15 {
		return "Hey there! I'm muh·ROH, your Marault Intelligence assistant. I'm here to help you explore our services and answer questions. What can I help you with today?"
	}

	// INAPPROPRIATE / HATEFUL / PERVERTED CONTENT HANDLER
	inappropriatePatterns := []string{"hate", "racist", "sexist", "slur", "lewd", "perverted", "nsfw", "adult", "porn", "sexual", "harass", "abuse", "violence", "kill", "rape", "assault", "discriminat", "offensive", "degrad", "retard"}
	if contains(userLower, inappropriatePatterns) {
		return "I don't feel comfortable discussing that topic. However, I'm here to help with anything related to Marault Intelligence and our business services. Please visit our services page to see how we can assist you."
	}

	// CRITICISM / INSULTS / SKEPTICISM HANDLER
	if contains(userLower, []string{"you suck", "terrible", "awful", "horrible", "garbage", "scam", "fraud", "waste", "stupid", "dumb", "useless", "you're wrong", "you're bad", "pathetic", "mediocre"}) ||
		(contains(userLower, []string{"doubt", "skeptic", "skeptical", "don't believe", "not convinced", "doesn't work"}) && len(userMessage) > 10) {
		return "We appreciate your perspective. If you have specific concerns about our approach or services, we'd love to discuss them. Visit inquire page to connect with our team and we can address any questions you might have."
	}

	// PRICING & COST QUESTIONS
	if contains(userLower, []string{"cost", "price", "pricing", "charge", "fee", "budget", "investment", "how much"}) {
		return "Our pricing varies based on the scope and complexity of your project. We offer flexible engagement models for startups, mid-market companies, and enterprises. To discuss pricing tailored to your needs, please fill out our inquiry form at inquire page and our team will provide a custom quote."
	}

	// TIMELINE & IMPLEMENTATION QUESTIONS
	if contains(userLower, []string{"timeline", "how long", "duration", "weeks", "months", "implement", "launch", "delivery"}) {
		return "Every project is unique based on scope and complexity. Let's discuss your specific needs and timeline at the inquire page so we can give you an accurate estimate."
	}

	// DATA SECURITY & PRIVACY
	if contains(userLower, []string{"security", "privacy", "data protection", "gdpr", "compliance", "safe", "confidential"}) {
		return "Data security is paramount to us. We follow industry best practices including encrypted data transmission, secure storage protocols, and strict confidentiality agreements. We're compliant with major regulations including GDPR and SOC 2. For detailed security documentation, please reach out via contact page."
	}

	// INTEGRATION & TECHNICAL
	if contains(userLower, []string{"integration", "api", "integrate", "connect", "system", "software", "platform", "salesforce", "tableau", "power bi"}) {
		return "We integrate with most major platforms including Salesforce, Tableau, Power BI, Google Analytics, and custom databases. Our team can connect your existing systems to create unified dashboards and workflows. Share your tech stack at inquire page so we can assess integration needs."
	}

	// INDUSTRY-SPECIFIC QUESTIONS
	if contains(userLower, []string{"industry", "finance", "ecommerce", "retail", "healthcare", "saas", "real estate", "professional services"}) {
		return "We've worked across diverse industries including Finance, E-Commerce, Retail, Healthcare, SaaS, Real Estate, and Professional Services. Each industry has unique data challenges—we customize solutions to address your specific sector. Tell us about your industry at inquire page."
	}

	// DATA VISIBILITY AUDIT - WHAT IS IT?
	if contains(userLower, []string{"what is data visibility audit", "data visibility audit", "what's a data audit", "data audit definition", "tell me about data audit", "understand data visibility", "data audit meaning", "how does audit work", "what does audit do", "audit my data"}) {
		return "A Data Visibility Audit is a structured review of how your company tracks, defines, and uses data to make decisions. We examine what data exists, where it lives, how it flows, how reliable it is, and whether leadership can trust it. Most companies discover their reporting is fragmented, inconsistent, or unclear. This engagement fixes that and establishes one trusted foundation for decisions. In simple terms: we make your numbers make sense. At Marault Intelligence, we specialize in transforming chaotic data landscapes into clear, decision-ready systems."
	}

	// DATA VISIBILITY AUDIT - WHY NEEDED?
	if contains(userLower, []string{"why do i need", "why would a company need", "do we need a data audit", "need data visibility", "why is this important", "should we audit data", "why audit", "benefits of audit", "do i need this", "is this necessary"}) {
		return "Most companies reach a point where reports don't match, teams argue over numbers, dashboards exist but don't help decisions, spreadsheets dominate workflows, and leadership hesitates to act. Research shows poor data quality directly impacts decision-making and productivity. At Marault Intelligence, we help companies eliminate this chaos by addressing the root causes—definitions, consistency, structure, and governance. We've worked with companies across Finance, SaaS, E-Commerce, and Professional Services to transform their data foundations. Does this sound familiar to your organization?"
	}

	// DATA VISIBILITY - WHAT DOES IT MEAN?
	if contains(userLower, []string{"what does data visibility mean", "data visibility meaning", "visibility definition"}) {
		return "Data visibility means leadership can see performance clearly, understand what changed, trust the numbers, and act quickly. Without visibility, decisions slow down and confidence drops. It's the difference between having data and actually using it to drive decisions."
	}

	// DATA VISIBILITY AUDIT - WHAT HAPPENS DURING?
	if contains(userLower, []string{"what happens during", "audit process", "what do you do", "how does the audit work", "engagement process"}) {
		return "During the audit we: 1) Map your data ecosystem, 2) Review reporting workflows, 3) Define decision-critical KPIs, 4) Evaluate data quality and consistency, 5) Identify structural issues, 6) Build a clean executive reporting layer, 7) Deliver a stabilization roadmap. Reach out at the inquire page to discuss the details of your project."
	}

	// DATA VISIBILITY AUDIT - EXECUTIVE BENEFITS
	if contains(userLower, []string{"what will leadership get", "executive benefits", "what does leadership receive", "leadership outcome"}) {
		return "Leadership receives: one trusted dashboard, clear KPI definitions, an insight summary showing what matters now, visibility into risk areas, and a prioritized fix roadmap. The goal is faster, more confident decision-making backed by reliable data."
	}

	// AUDIT VS DASHBOARD
	if contains(userLower, []string{"difference between audit and dashboard", "why audit before dashboard", "dashboard vs audit", "why not just build dashboard"}) {
		return "Most dashboards fail because the foundation is wrong. Common issues include inconsistent metric definitions, broken joins, duplicated logic, unclear ownership, and low data quality. The audit fixes these foundation issues first—then builds reporting that actually lasts. It's the difference between a quick fix and a reliable system."
	}

	// AUDIT - WHAT PROBLEMS DOES IT SOLVE?
	if contains(userLower, []string{"problems the audit solves", "what problems", "solves what", "inconsistent reporting", "teams debating numbers"}) {
		return "The audit solves: inconsistent reporting, slow decision-making, lack of alignment across teams, unreliable dashboards, and unclear priorities. If your teams debate which numbers are correct or leadership doesn't trust the data, the audit addresses these exact problems."
	}

	// AUDIT - WHO IS IT BEST FOR?
	if contains(userLower, []string{"who is this for", "best fit", "right for us", "is this for my company"}) {
		return "Best fit: growing companies, leadership teams lacking clarity, companies scaling operations, and firms with fragmented systems. It's ideal if you're noticing reporting inconsistencies as you scale. It's less necessary if you already have strong governance and clean infrastructure in place."
	}

	// AUDIT - TECHNICAL EVALUATION
	if contains(userLower, []string{"what do you evaluate technically", "technical evaluation", "what systems do you review"}) {
		return "We review: data models, table structures, join logic, ETL workflows, reporting layers, metric definitions, and governance processes. We assess both your infrastructure and reporting outputs to understand the full picture of data quality."
	}

	// DATA LINEAGE QUESTIONS
	if contains(userLower, []string{"data lineage", "what is lineage", "lineage mapping"}) {
		return "Data lineage shows how data moves from source through transformation to reporting. It matters because it prevents reporting errors, helps debug issues faster, ensures transparency, and supports governance. Modern data teams treat lineage as essential infrastructure for reliability."
	}

	// DATA QUALITY EVALUATION
	if contains(userLower, []string{"data quality", "quality evaluation", "how do you evaluate quality"}) {
		return "We assess completeness (missing fields), consistency (IDs, formats), accuracy, duplication, and rule alignment. These dimensions align with standard data quality frameworks used across enterprise analytics. We prioritize issues by business impact."
	}

	// KPI STANDARDIZATION
	if contains(userLower, []string{"kpi standardization", "standardize metrics", "consistent metrics", "metric definitions"}) {
		return "KPI standardization involves defining calculation logic, aligning metrics across teams, documenting assumptions, and establishing governance. This is critical because inconsistent definitions are one of the most common sources of reporting conflict. Teams stop debating and start trusting the numbers."
	}

	// INFRASTRUCTURE CHANGES
	if contains(userLower, []string{"do you change infrastructure", "change our stack", "modify systems", "do you modify our tools"}) {
		return "Usually no. We improve clarity using your existing tools wherever possible. We assess your pipelines, storage, transformation logic, and reporting outputs—but focus on optimizing what you have rather than requiring new infrastructure."
	}

	// AUDIT TIMELINE DETAILED
	if contains(userLower, []string{"audit timeline", "how long audit", "audit duration", "audit timeframe"}) {
		return "The audit typically follows phases: discovery (understand your systems), evaluation (analyze quality and structure), cleanup (fix critical issues), reporting build (create the dashboard), and stabilization plan (handoff). For exact timelines and details, visit the inquire page to connect with our team."
	}

	// TEAM INVESTMENT
	if contains(userLower, []string{"how much time", "team investment", "time commitment", "what do we need to do"}) {
		return "Minimal time investment from your team. We handle most analysis independently and schedule focused sessions for alignment and context. Your main involvement is 1-2 hours total during discovery and review phases."
	}

	// SYSTEMS COMPATIBILITY
	if contains(userLower, []string{"what systems do you work with", "compatible systems", "crm", "data warehouse", "BI tool", "tableau", "power bi"}) {
		return "We work with most systems: CRMs, data warehouses, BI tools (Tableau, Power BI, Looker), spreadsheets, and internal systems. We're agnostic to your tech stack—the principles of data visibility and quality apply everywhere."
	}

	// AUDIT COST/PRICING
	if contains(userLower, []string{"audit cost", "audit price", "how much audit", "audit pricing"}) {
		return "Audit pricing depends on complexity and scope. We tailor pricing based on: systems involved, data maturity, reporting complexity, and volume. Rather than one-size-fits-all pricing, we provide custom quotes after discovery. Start a conversation at inquire page."
	}

	// WHY NOT DO INTERNALLY?
	if contains(userLower, []string{"why not solve internally", "do it ourselves", "internal team", "in-house"}) {
		return "Internal teams often lack bandwidth, lack objectivity (they're embedded in the problem), or have competing priorities. We provide speed, fresh perspective, and structured methodology. Plus, an external audit carries credibility with leadership and removes politics from data decisions."
	}

	// IS THIS JUST CONSULTING?
	if contains(userLower, []string{"is this just consulting", "just consulting", "do you just advise", "operational deliverables"}) {
		return "No. You leave with operational systems, working dashboards, governance structure, and a roadmap. This is about building lasting infrastructure—not just recommendations. Leadership gets immediate visibility through the new dashboard, and teams get clarity through standardized definitions."
	}

	// AFTER THE AUDIT
	if contains(userLower, []string{"after the audit", "next steps", "what happens next", "post-audit"}) {
		return "After the audit, clients typically implement the fix roadmap, expand reporting to new areas, build forecasting systems, and scale analytics maturity. The audit becomes your foundation for deeper analytics work. Many clients contract for ongoing support or expand into our other services."
	}

	// FIRST STEP
	if contains(userLower, []string{"first step", "how do we start", "getting started", "next move"}) {
		return "The first step is a short discovery conversation. We'll understand your biggest challenges, current tools, reporting gaps, and goals. This takes 15-20 minutes and helps us scope exactly what you need. Ready to talk? Visit inquire page or contact page."
	}

	// GENERIC DATA AUDIT VISIBILITY
	if contains(userLower, []string{"data audit", "visibility", "understand data", "what data", "data inventory"}) {
		return "Our Data Visibility Audit analyzes your current data infrastructure to identify what data you have, where it lives, how it flows, how reliable it is, and whether leadership can trust it. We deliver actionable recommendations along with a working dashboard. Visit the inquire page to discuss your specific needs and get a customized timeline."
	}

	// ==========================================
	// EXECUTIVE DASHBOARD QUESTIONS
	// ==========================================

	// EXECUTIVE DASHBOARD - WHAT IS IT?
	if contains(userLower, []string{"what is an executive dashboard", "executive dashboard definition", "what's an executive dashboard", "executive reporting", "what's a dashboard", "dashboard meaning", "executive report", "tell me about dashboards", "understand dashboard"}) {
		return "An executive dashboard is a simplified reporting system designed for leadership decision-making. Unlike standard dashboards, it focuses on the few metrics that actually matter, what changed recently, why it changed, and what action leadership should take. Research shows executives prefer high-level summaries and trends over operational detail. At Marault Intelligence, we design executive dashboards that prioritize clarity and speed over volume. Our dashboards become the core operating view your leadership team relies on for strategy and decisions."
	}

	// EXECUTIVE DASHBOARD - VS NORMAL DASHBOARD
	if contains(userLower, []string{"executive dashboard vs normal", "difference between dashboards", "how is executive different", "vs regular dashboard", "why different", "not just a dashboard", "compared to normal"}) {
		return "Most dashboards track too many metrics, lack context, overwhelm users, and don't guide decisions. Executive dashboards prioritize signal over noise, align metrics across teams, surface drivers and risks, and support fast decisions. The difference is focus: operational dashboards inform, executive dashboards decide. At Marault Intelligence, we specialize in building dashboards that your leadership team actually uses—not ones that sit unused. Clients see significant improvements in decision speed and team alignment after implementation."
	}

	// EXECUTIVE DASHBOARD - WHY NEEDED?
	if contains(userLower, []string{"why executive dashboard", "why would we need this", "when do we need executive dashboard"}) {
		return "Companies usually come to us when leadership spends too long reviewing reports, numbers are inconsistent, meetings focus on data disputes instead of decisions, or reporting feels reactive. This service creates structure, alignment, and clarity. If your leadership meetings drag on because of reporting confusion, you need this."
	}

	// EXECUTIVE DASHBOARD - LEADERSHIP GAINS
	if contains(userLower, []string{"what will leadership gain", "leadership benefits", "what does leadership get", "executive benefits"}) {
		return "Leadership gains a fast decision system, shared definitions across teams, clear performance visibility, early warning signals, and confidence in reporting. The goal is faster alignment and stronger decisions. Most clients tell us their leadership meetings become 50% shorter and more decision-focused."
	}

	// DECISION-READY REPORTING
	if contains(userLower, []string{"decision-ready reporting", "what does decision ready mean", "decision ready"}) {
		return "Decision-ready reporting means clear trends, defined thresholds, context around changes, and actionable insight. It removes interpretation friction so leadership can act immediately without spending time analyzing the data. Every metric tells a story: what changed, why, and what to do next."
	}

	// EXECUTIVE DASHBOARD - USAGE FREQUENCY
	if contains(userLower, []string{"how often used", "usage frequency", "how often will leadership use"}) {
		return "Most clients use executive dashboards weekly, monthly, or quarterly depending on their business cycle. It typically becomes the core operating view leadership relies on. Some clients check it daily during critical periods, others during scheduled reviews. The point is it becomes their standard view, not an ad-hoc report."
	}

	// EXECUTIVE DASHBOARD - MEETINGS
	if contains(userLower, []string{"help leadership meetings", "leadership meetings", "board meetings", "executive meetings"}) {
		return "Executive dashboards reduce time spent reviewing data, confusion across teams, and reactive discussions. They increase alignment, clarity, and speed. Instead of 30 minutes debating what the numbers mean, leadership spends 5 minutes reviewing the dashboard and 25 minutes discussing strategy. It transforms the nature of leadership meetings."
	}

	// EXECUTIVE DASHBOARD - METRICS
	if contains(userLower, []string{"what metrics", "which metrics do you include", "how many metrics", "kpi count"}) {
		return "We define a small KPI set based on strategic goals, operating drivers, and accountability structure. Research shows effective executive dashboards typically contain 5-12 metrics, not 50+. This prevents overload while ensuring leadership sees what matters. The art is knowing what to exclude."
	}

	// NARRATIVE REPORTING LAYER
	if contains(userLower, []string{"narrative reporting", "narrative layer", "what is narrative reporting"}) {
		return "The narrative layer explains what changed, why it changed, and what matters next. Instead of just showing numbers, we provide context. For example: 'Revenue up 12% (seasonal), churn down to 2.1% (new support process), customer acquisition cost stable.' This reflects how executives actually think about performance."
	}

	// EXECUTIVE DASHBOARD - REPLACE INTERNAL?
	if contains(userLower, []string{"replace internal reporting", "will this replace", "other reporting", "team dashboards"}) {
		return "No. We separate executive signal from team analysis. Both remain useful. Leadership gets the simplified executive view while teams keep their operational dashboards. This prevents information fragmentation—everyone sees their level of detail."
	}

	// EXECUTIVE DASHBOARD - TECHNICAL DESIGN
	if contains(userLower, []string{"how do you design dashboards", "dashboard design", "technical design", "design principles"}) {
		return "We follow proven BI design principles: clear visual hierarchy, minimal cognitive load, strong contrast and readability, consistent layouts, and trend emphasis. These align with best practices from BI research and data visualization standards. Good design gets out of the way and lets insights shine."
	}

	// EXECUTIVE DASHBOARD - TOOLS
	if contains(userLower, []string{"what tools do you use", "dashboard tools", "BI platforms", "tableau power bi"}) {
		return "We typically work within your stack—BI platforms (Tableau, Power BI, Looker), data warehouses, spreadsheets, or internal tools. We prioritize integration over replacement. If you already have infrastructure, we build on it rather than require new tools."
	}

	// EXECUTIVE DASHBOARD - DATA TRUST
	if contains(userLower, []string{"ensure trust", "data trust", "trust the data", "validate data"}) {
		return "We implement metric governance, validation checks, definition control, and structured review cadence. Trust comes from consistency and transparency. Every metric has a clear owner, defined calculation, and documented assumptions. This aligns with established governance frameworks."
	}

	// DRILL-PATHS
	if contains(userLower, []string{"drill-paths", "drill path", "drill down", "explore details"}) {
		return "Drill-paths allow teams to explore details, validate trends, and identify drivers while leadership stays focused on signal. For example, a CEO sees 'revenue up 12%' but can drill down to see which products or regions are driving growth. It's the bridge between simplicity and depth."
	}

	// METRIC DRIFT PREVENTION
	if contains(userLower, []string{"metric drift", "prevent drift", "definition drift", "metric consistency"}) {
		return "We enforce clear definitions, ownership, documentation, and governance cadence. Metric drift is a common failure point in BI systems—what revenue means shifts over time, teams calculate differently, and reporting loses credibility. Governance prevents this."
	}

	// EXECUTIVE DASHBOARD - ENGAGEMENT PROCESS
	if contains(userLower, []string{"engagement process", "what does the engagement look like", "how does engagement work"}) {
		return "Typical phases: 1) KPI alignment with leadership, 2) Design and structure, 3) Build, 4) Narrative setup, 5) Rollout, 6) Cadence support. We collaborate closely but handle most build work independently. Your main time investment is in alignment sessions."
	}
	// EXECUTIVE DASHBOARD - TIMELINE
	if contains(userLower, []string{"dashboard timeline", "how long dashboard", "dashboard duration"}) {
		return "Every project is unique based on complexity and systems. We define timelines after our initial discovery conversation. For exact details, visit the inquire page to discuss your specific needs with our team."
	}

	// EXECUTIVE DASHBOARD - TEAM EFFORT
	if contains(userLower, []string{"team effort required", "what's required from team", "time commitment dashboard"}) {
		return "Mostly alignment sessions and feedback. We handle most build work independently. Your team's main responsibility is defining what metrics matter and validating the final product. Typically 2-3 hours total from your leadership team over the engagement."
	}

	// EXECUTIVE DASHBOARD - BUILD INTERNALLY
	if contains(userLower, []string{"build internally", "why not build ourselves", "build in-house"}) {
		return "Internal teams often lack bandwidth, neutrality, or design specialization. We bring speed, structure, and perspective. Plus, an external view helps depoliticize metrics—it's easier for teams to accept definitions from us than from an internal stakeholder."
	}

	// EXECUTIVE DASHBOARD - JUST DESIGN?
	if contains(userLower, []string{"just a design project", "is this design", "design only"}) {
		return "No. It's an operating system for decision-making. We design, build, validate data integrity, establish governance, and set up the cadence for ongoing use. It's not about making pretty charts—it's about creating lasting structure for how leadership makes decisions."
	}

	// EXECUTIVE DASHBOARD - COST
	if contains(userLower, []string{"executive dashboard cost", "dashboard pricing", "dashboard price"}) {
		return "Pricing depends on scope, systems, and complexity. We tailor engagements based on your needs and infrastructure. Rather than quote a fixed price, we discuss your goals first and scope accordingly. Chat with us at /inquire for a custom quote."
	}

	// EXECUTIVE DASHBOARD - DO YOU NEED IT?
	if contains(userLower, []string{"do we need this", "is this for us", "do we need executive dashboard"}) {
		return "Common signals you need this: long reporting meetings, inconsistent numbers, reactive leadership decisions, or teams spending more time debating data than discussing strategy. If your leadership meetings feel like data reviews instead of decision forums, you probably need this."
	}

	// EXECUTIVE DASHBOARD - AFTER LAUNCH
	if contains(userLower, []string{"after launch", "post-launch", "what happens after"}) {
		return "Clients often continue with cadence support (we attend and facilitate reviews), expand into forecasting, or move into deeper analytics. The dashboard becomes your foundation for more sophisticated analytics. Many clients scale from here into our other services like revenue analytics or decision modeling."
	}

	// ==========================================
	// FORECASTING & DECISION MODELING QUESTIONS
	// ==========================================

	// FORECASTING - WHAT IS IT?
	if contains(userLower, []string{"what is forecasting", "decision modeling", "what is forecasting decision modeling"}) {
		return "Forecasting & Decision Modeling is a system that helps you predict future outcomes and make better decisions under uncertainty. Instead of guessing or reacting, it shows what is likely to happen, what could happen, and what decisions you should make. It turns planning into a structured process instead of intuition."
	}

	// FORECASTING - WHY NOT JUST A NUMBER?
	if contains(userLower, []string{"single forecast", "just one number", "single prediction", "why ranges not numbers"}) {
		return "A single number is almost always misleading. Real forecasting focuses on ranges (best case and worst case), assumptions, and uncertainty. Modern forecasting best practices emphasize scenario-based planning because business environments are dynamic."
	}

	// FORECASTING - WHAT DOES DECISION MODELING MEAN?
	if contains(userLower, []string{"decision modeling", "what is decision modeling", "decision model"}) {
		return "Decision modeling connects choices, constraints, and outcomes. It helps answer questions like: Should we hire now or wait? Should we increase pricing? Can we support more demand? It's the framework that ties decisions to their consequences."
	}

	// FORECASTING - WHAT WILL LEADERSHIP GAIN?
	if contains(userLower, []string{"forecast leadership benefits", "what does leadership gain", "decision modeling benefits"}) {
		return "Leadership gains clarity on future outcomes, structured decision-making, early visibility into risk, and confidence in planning. Instead of reacting, leadership can plan ahead."
	}

	// FORECASTING - HOW DOES IT IMPROVE DECISIONS?
	if contains(userLower, []string{"improve decision making", "better decisions", "how does modeling help decisions"}) {
		return "It provides clear trade-offs, quantified scenarios, decision thresholds, and visibility into impact. This aligns with modern FP&A practices where decisions are tied to modeled outcomes."
	}

	// FORECASTING - WHAT DECISIONS DOES IT HELP WITH?
	if contains(userLower, []string{"hiring decision", "pricing strategy", "expansion", "inventory planning", "capacity planning", "what decisions"}) {
		return "Examples include hiring and headcount planning, pricing strategy, expansion into new markets, inventory or capacity planning, and marketing spend allocation. Any decision involving uncertainty and future planning benefits from this approach."
	}

	// FORECASTING - WHAT IS SCENARIO PLANNING?
	if contains(userLower, []string{"scenario planning", "what is scenario", "scenarios"}) {
		return "Scenario planning models different possible futures: base case, upside, and downside. Each scenario includes drivers, assumptions, outcomes, and actions. This approach is widely used in strategic planning and risk management."
	}

	// FORECASTING - WHAT DOES A MODEL INCLUDE?
	if contains(userLower, []string{"what does forecasting model include", "model components", "what is in a model"}) {
		return "A good model includes drivers (inputs that influence outcomes), assumptions, constraints, outputs, and scenario logic. Each component works together to create a realistic forecast."
	}

	// FORECASTING - WHAT ARE DRIVERS?
	if contains(userLower, []string{"what are drivers", "drivers forecast", "input variables"}) {
		return "Drivers are the variables that control outcomes. Examples include demand, conversion rates, pricing, churn, and capacity. Driver-based modeling is a standard approach in financial and operational forecasting."
	}

	// FORECASTING - WHAT IS DECISION THRESHOLD?
	if contains(userLower, []string{"decision threshold", "threshold", "trigger point"}) {
		return "A decision threshold is a trigger point. For example: If demand drops below X, reduce spend. If CAC exceeds Y, adjust strategy. This creates clear action rules instead of ambiguity."
	}

	// FORECASTING - EARLY RISK INDICATORS
	if contains(userLower, []string{"early risk", "early indicators", "leading indicators", "risk signals"}) {
		return "Early indicators signal problems before they show up in results. Examples include declining conversion, slowing demand, and increased churn signals. These are often called leading indicators in analytics."
	}

	// FORECASTING - HOW DO YOU BUILD MODELS?
	if contains(userLower, []string{"how do you build models", "build forecasting models", "model building process"}) {
		return "We use driver-based models, scenario analysis, sensitivity testing, and constraint modeling. These are standard techniques in FP&A and operational planning."
	}

	// FORECASTING - WHAT IS SENSITIVITY ANALYSIS?
	if contains(userLower, []string{"sensitivity analysis", "sensitivity testing", "what if analysis"}) {
		return "Sensitivity analysis tests which variables impact outcomes most. For example: if pricing changes, what happens to revenue? If demand drops, what happens to capacity? This helps identify the biggest risks and opportunities."
	}

	// FORECASTING - HOW DO YOU HANDLE UNCERTAINTY?
	if contains(userLower, []string{"handle uncertainty", "uncertainty modeling", "confidence levels"}) {
		return "We use ranges instead of point estimates, model multiple scenarios, and define confidence levels. This aligns with best practices in forecasting and risk modeling."
	}

	// FORECASTING - WHAT MAKES A MODEL GOOD?
	if contains(userLower, []string{"what makes model good", "good model", "model quality"}) {
		return "A good model is explainable, maintainable, aligned to the business, and adaptable. Not overly complex or fragile. It serves the business, not the other way around."
	}

	// FORECASTING - WHAT'S WRONG WITH MOST MODELS?
	if contains(userLower, []string{"problems with forecasting", "why forecasts fail", "forecasting mistakes"}) {
		return "Common issues include overly complex models, hard to maintain systems, disconnected from operations, and based on unrealistic assumptions. Many forecasts fail because they're black boxes that don't connect to actual decision-making."
	}

	// FORECASTING - ENGAGEMENT PROCESS
	if contains(userLower, []string{"forecasting engagement", "forecasting process", "how do we work"}) {
		return "Typical phases: 1) Define drivers and assumptions, 2) Build model structure, 3) Create scenarios, 4) Define decision thresholds, 5) Implement planning tools, 6) Establish cadence. We work closely with your team throughout."
	}

	// FORECASTING - TIMELINE
	if contains(userLower, []string{"forecasting timeline", "how long model", "forecasting duration"}) {
		return "Every project is unique based on complexity and scope. We define timelines after understanding your specific needs. Visit the inquire page to discuss your requirements and get an accurate estimate."
	}

	// FORECASTING - DATA REQUIREMENTS
	if contains(userLower, []string{"what data do you need", "data requirements forecasting", "historical data"}) {
		return "Typically we need historical performance data, operational data, financial data, and your assumptions. We work with what you have and help you identify gaps."
	}

	// FORECASTING - EXCEL VS FORECASTING SYSTEM
	if contains(userLower, []string{"excel forecasting", "why not excel", "spreadsheet vs system"}) {
		return "Basic spreadsheets don't handle uncertainty well, don't connect decisions to outcomes, and break easily. Our system is structured, scalable, and decision-focused."
	}

	// FORECASTING - BUILD INTERNALLY?
	if contains(userLower, []string{"build forecasting internally", "build in house", "internal forecasting"}) {
		return "Internal teams often focus on reporting rather than modeling, lack time to build robust systems, and don't formalize decision frameworks. We bring specialized expertise and dedicated resources."
	}

	// FORECASTING - FOR LARGE COMPANIES ONLY?
	if contains(userLower, []string{"large companies only", "for big companies", "company size requirements"}) {
		return "No. This is especially valuable for growing companies, companies making frequent decisions, and companies managing uncertainty. Size doesn't matter—decision complexity does."
	}

	// FORECASTING - DO YOU NEED IT?
	if contains(userLower, []string{"do we need forecasting", "planning feels reactive", "forecasts unreliable"}) {
		return "Common signals you need this: planning feels reactive, forecasts are unreliable, decisions are debated without clarity, or growth and capacity are unpredictable. If you're making decisions without a clear framework, this helps."
	}

	// FORECASTING - AFTER IMPLEMENTATION
	if contains(userLower, []string{"after forecasting", "forecasting next steps", "post implementation"}) {
		return "Clients typically improve forecasting accuracy, make faster decisions, align leadership, and expand into deeper analytics. This service often becomes the foundation for more sophisticated planning and strategy work."
	}

	// REVENUE & CUSTOMER ANALYTICS
	if contains(userLower, []string{"revenue", "customer analytics", "churn", "customer behavior", "sales", "growth", "retention"}) {
		return "Our Revenue & Customer Analytics helps you understand revenue drivers, customer lifetime value, churn patterns, and growth opportunities. We create dashboards that show real-time metrics and predictive insights. This helps you make data-driven decisions to increase profitability. Visit services page on revenue customer analytics."
	}

	// EXECUTIVE DASHBOARDS & REPORTING
	if contains(userLower, []string{"executive", "dashboard", "reporting", "metrics", "kpi", "real-time", "board", "c-suite"}) {
		return "Executive Dashboards & Reporting gives C-suite leaders the insights they need in real-time. We build custom dashboards with your key metrics, automated reports, and predictive insights. Perfect for board meetings and strategic planning. Check out services page on executive dashboards reporting."
	}

	// FORECASTING & MODELING
	if contains(userLower, []string{"forecasting", "modeling", "predict", "scenario", "forecast", "planning", "budget", "demand"}) {
		return "Our Forecasting & Decision Modeling services use advanced analytics to predict future trends, model scenarios, and optimize planning. Whether it's revenue forecasts, demand planning, or strategic scenarios, we help you make confident decisions. Learn more at services page on forecasting decision modeling."
	}

	// PRIVATE CLIENT ANALYTICS - BEGINNER QUESTIONS
	if contains(userLower, []string{"what is private client analytics", "private client analytics definition", "what do you mean by private client"}) {
		return "Private Client Analytics is a structured way to organize and interpret a client's full financial picture. It brings together accounts, holdings, cash flows, risk exposures, and performance into one clear view so decisions are based on facts instead of fragments."
	}

	if contains(userLower, []string{"who is private client analytics for", "is this for me", "private client who is it for"}) {
		return "This is for private clients who have meaningful financial complexity: multiple accounts, multiple institutions, concentrated positions, liquidity decisions, tax considerations, or long-term planning questions. It's especially useful when someone has plenty of information but no single, decision-ready view."
	}

	if contains(userLower, []string{"what problem does private client solve", "problem private client", "why would i need private client analytics"}) {
		return "Many private clients get statements, advisor commentary, account portals, and market updates, but still cannot answer basic decision questions quickly like: What is actually driving my performance? Where am I overexposed? How much risk am I really taking? This service solves that by turning scattered information into a coherent framework."
	}

	if contains(userLower, []string{"is private client investment management", "investment management vs private client", "does private client manage money"}) {
		return "Not exactly. Private Client Analytics is best described as an analytics and decision-support layer. It helps a client see, evaluate, and structure decisions more clearly. It can complement a financial advisor, wealth manager, CPA, attorney, or family office rather than replace them."
	}

	if contains(userLower, []string{"why not broker dashboard", "broker portal vs private client", "isn't broker enough"}) {
		return "Most broker dashboards are account views, not decision systems. They often show balances and returns, but not the full cross-account picture, risk concentration, scenario impacts, or decision trade-offs."
	}

	// PRIVATE CLIENT ANALYTICS - PERFORMANCE QUESTIONS
	if contains(userLower, []string{"consolidated performance view", "what is consolidated performance", "cross account performance"}) {
		return "It means one coherent picture across accounts, entities, and holdings. Instead of reviewing multiple statements separately, the client can see total exposure, performance, and key changes in one place."
	}

	if contains(userLower, []string{"performance attribution", "what is attribution", "what drove my returns"}) {
		return "Performance attribution explains what drove results. Instead of just seeing that a portfolio gained or lost value, attribution helps answer: Was it asset allocation? A single concentrated holding? Sector exposure? Cash drag? Timing? Currency or tax effects? That matters because better decisions come from understanding why results happened, not just what happened."
	}

	if contains(userLower, []string{"why benchmarks important", "benchmark question", "what is a benchmark"}) {
		return "A benchmark helps answer whether results are strong or weak relative to the client's objectives and risk profile. A return number without context can be misleading. Good benchmarking should be aligned to what the client is actually trying to achieve."
	}

	if contains(userLower, []string{"which parts helping hurting", "portfolio contribution", "which holdings are helping"}) {
		return "Yes. That is one of the core uses. The goal is to separate noise from true contribution and show which exposures, accounts, or positions are driving outcomes."
	}

	// PRIVATE CLIENT ANALYTICS - RISK QUESTIONS
	if contains(userLower, []string{"risk exposure mapping", "what is exposure mapping", "how do you map risk"}) {
		return "Risk exposure mapping shows where the portfolio is concentrated and how different holdings may behave together. That includes exposure by asset class, sector, geography, strategy, account type, and sometimes by underlying factor or theme."
	}

	if contains(userLower, []string{"why concentration matters", "concentration risk", "concentrated position risk"}) {
		return "Because a portfolio that is heavily exposed to one stock, sector, or theme can suffer outsized losses if that area underperforms. A large portion of assets in one stock represents concentration risk, and diversification is a core way to manage investment risk."
	}

	if contains(userLower, []string{"hidden correlation", "correlation risk", "investments moving together"}) {
		return "Two investments may look different on paper but still tend to move together under stress. Hidden correlation means a client may think they are diversified when, in practice, several holdings are exposed to the same underlying risk. This is important because diversification depends heavily on how assets move relative to each other."
	}

	if contains(userLower, []string{"how think about diversification", "diversification definition", "what is diversification"}) {
		return "Diversification means spreading investments among and within asset classes to help reduce risk. It's a regular adjustment process to keep allocation aligned over time."
	}

	if contains(userLower, []string{"help rebalancing", "rebalancing decisions", "when should i rebalance"}) {
		return "Yes, in an analytic sense. We can frame when a client may want to review exposures, trim concentration, or compare alternative uses of capital. The point is not constant trading; it is disciplined decision-making."
	}

	// PRIVATE CLIENT ANALYTICS - SCENARIO QUESTIONS
	if contains(userLower, []string{"scenario analysis", "what is scenario analysis", "how scenarios work"}) {
		return "Scenario analysis shows how a portfolio or financial plan may behave under different conditions, such as: an equity drawdown, interest-rate changes, weaker liquidity, a concentrated position falling sharply, or a large near-term cash need. It helps a client see possible outcomes before they happen."
	}

	if contains(userLower, []string{"stress testing", "what is stress test", "stress test portfolio"}) {
		return "Stress testing is a more severe form of scenario analysis. It asks how the portfolio performs under adverse conditions and what vulnerabilities appear. It's important when correlations change under pressure."
	}

	if contains(userLower, []string{"why scenarios matter private", "why stress testing important", "scenario importance"}) {
		return "Because high-net-worth decisions are often path-dependent. A client may be able to tolerate volatility in theory, but not at the exact moment they need liquidity, have tax consequences, or are making a major allocation choice."
	}

	if contains(userLower, []string{"what if modeling", "can you model what if", "what if analysis"}) {
		return "Yes. For example: What if I trim a concentrated holding? What if I need more cash over the next 12 months? What if markets fall 20%? What if I change my allocation? What if I delay or accelerate a major decision?"
	}

	// PRIVATE CLIENT ANALYTICS - ALLOCATION QUESTIONS
	if contains(userLower, []string{"allocation framework", "trade off framework", "what is allocation framework"}) {
		return "It is a structured way to compare options. Instead of asking only Which investment could go up more?, it asks: What is the objective? What risk is being taken? What liquidity is needed? What tax effect matters? What is the opportunity cost?"
	}

	if contains(userLower, []string{"opportunity cost", "what is opportunity cost", "meaning opportunity cost"}) {
		return "Opportunity cost is the cost of choosing one use of capital over another. In private wealth decisions, that could mean comparing: staying concentrated vs diversifying, holding cash vs investing, selling now vs deferring, or prioritizing tax efficiency vs liquidity."
	}

	if contains(userLower, []string{"help decide hold reduce rebalance", "when to rebalance", "hold reduce decision"}) {
		return "Yes. The service is designed to create clearer decision rules and thresholds so the client is not relying only on intuition or headlines."
	}

	// PRIVATE CLIENT ANALYTICS - TAX QUESTIONS
	if contains(userLower, []string{"tax aware", "what does tax aware mean", "tax aware analytics"}) {
		return "It means decisions are evaluated based on after-tax outcomes, not just pre-tax returns. That matters because investment income and capital gains can be taxed differently."
	}

	if contains(userLower, []string{"why after tax matters", "after tax importance", "why tax outcomes matter"}) {
		return "Because the same investment return can lead to very different net results depending on account type, timing, gains, losses, dividends, and other tax factors. A client making large allocation decisions without considering taxes may misunderstand the true result."
	}

	if contains(userLower, []string{"help cash needs", "cash planning", "liquidity planning"}) {
		return "Yes. One of the most useful parts of this service is connecting portfolio decisions to cash-flow needs and timing. That helps avoid being forced into reactive decisions at the wrong time."
	}

	if contains(userLower, []string{"is this tax advice", "tax advice private client", "legal advice"}) {
		return "No. This should be framed as tax-aware analytics, not legal or tax advice. The goal is to surface the decision impact clearly and coordinate with the client's CPA or tax advisor where needed."
	}

	// PRIVATE CLIENT ANALYTICS - PROCESS QUESTIONS
	if contains(userLower, []string{"what does engagement look like", "engagement process", "how does this work"}) {
		return "A typical engagement includes: consolidating relevant financial data, organizing accounts and exposures, building a performance and risk view, developing scenario and stress-test frameworks, creating a clean recurring briefing format, and setting a monthly or quarterly review cadence."
	}

	if contains(userLower, []string{"what information do you need", "what data needed", "information required"}) {
		return "Usually: account statements or holdings exports, basic objective and liquidity context, entity or account structure if relevant, and any current benchmark, advisor reports, or planning materials."
	}

	if contains(userLower, []string{"how often updated", "update frequency", "when updated"}) {
		return "Often monthly or quarterly, depending on complexity and how active the decision cycle is."
	}

	if contains(userLower, []string{"what does client receive", "deliverables private client", "what will i get"}) {
		return "Typically: a consolidated view, performance attribution, risk exposure summary, scenario or stress-test outputs, a private briefing format, and a repeatable decision cadence."
	}

	// PRIVATE CLIENT ANALYTICS - COMPARISON QUESTIONS
	if contains(userLower, []string{"different from wealth manager", "vs wealth manager", "wealth manager comparison"}) {
		return "A wealth manager typically manages money, gives investment advice, or oversees planning relationships. Private Client Analytics is narrower and more decision-analytic: it gives the client a clearer, structured view of what is happening and what trade-offs exist."
	}

	if contains(userLower, []string{"different from family office", "vs family office", "family office comparison"}) {
		return "A family office can be broad and operational. This service is more focused: it builds a high-clarity analytics layer for visibility, performance interpretation, risk framing, and decision support."
	}

	if contains(userLower, []string{"why not statements quarterly", "statements vs private client", "quarterly review enough"}) {
		return "Because statements are fragmented and often backward-looking. This service is designed to improve coherence, speed, and decision quality."
	}

	if contains(userLower, []string{"only wealthy clients", "only ultra high net worth", "minimum wealth needed"}) {
		return "It is best suited for clients with real complexity, not necessarily only ultra-high-net-worth households. The key factor is whether the client has meaningful decisions, multiple moving parts, and a need for discretion and clarity."
	}

	// PRIVATE CLIENT ANALYTICS - TECHNICAL QUESTIONS
	if contains(userLower, []string{"evaluate concentration", "how to measure concentration", "concentration analysis"}) {
		return "We look at where risk is actually accumulating: individual positions, sectors, correlated exposures, illiquidity, and cross-account overlap."
	}

	if contains(userLower, []string{"correlation changes stress", "correlation stress testing", "correlation analysis"}) {
		return "That is often a critical part of the analysis. We specifically account for the importance of recalibrating assumed correlations in sensitivity analysis and stress testing."
	}

	if contains(userLower, []string{"benchmark relative objective relative", "benchmark vs objective", "performance vs objective"}) {
		return "Yes. Both can matter. A client may care about beating a benchmark, preserving capital, funding future needs, or reducing concentration risk. The analysis should reflect the actual objective."
	}

	if contains(userLower, []string{"entity level multi account", "multiple entities", "cross account analysis"}) {
		return "Yes. The service is well suited for clients with accounts across institutions, entities, or different wrappers where the total picture is otherwise hard to interpret."
	}

	// PRIVATE CLIENT ANALYTICS - OBJECTIONS
	if contains(userLower, []string{"already have advisor", "have wealth manager", "have advisor already"}) {
		return "That's completely fine. Private Client Analytics can complement an advisor by creating a cleaner decision-support layer across accounts, exposures, and scenarios."
	}

	if contains(userLower, []string{"just want understand risk", "understand risk better", "risk focused"}) {
		return "That's one of the clearest use cases. We can help surface concentration, hidden correlation, downside sensitivity, and the trade-offs behind possible adjustments."
	}

	if contains(userLower, []string{"multiple accounts no view", "accounts spread across", "no single view"}) {
		return "That is exactly the kind of complexity this service is built for. The first step is consolidating the picture so decisions stop relying on fragments."
	}

	if contains(userLower, []string{"don't need constant monitoring", "no ongoing monitoring", "don't want monitoring"}) {
		return "That aligns well with the service. The goal is not endless monitoring—it is a clean briefing view and a repeatable cadence that supports confident action."
	}

	// CUSTOM WEBSITE BUILD - BEGINNER QUESTIONS
	if contains(userLower, []string{"what is custom website build", "custom website definition", "what is a custom build"}) {
		return "A custom website build is a site designed and developed specifically around your business, offer, audience, and conversion goals rather than starting from a generic template. The point is not just appearance—it is to make your positioning clearer, your credibility stronger, and your path to inquiry easier."
	}

	if contains(userLower, []string{"why choose custom over template", "custom vs template", "template vs custom"}) {
		return "Templates are faster, but they often force your business into a pre-made structure. A custom site lets the structure, messaging, hierarchy, and visual system reflect how your business actually sells. That matters because users tend to scan pages rather than read them line by line, so structure and clarity have a major effect on whether they understand and trust the offer."
	}

	if contains(userLower, []string{"what problem custom website solve", "problem custom website", "why would i need custom"}) {
		return "Most sites are fine visually but weak strategically. They often leave visitors asking: What does this company actually do? Is this high quality? Is this for someone like me? What should I do next? A strong custom build solves those questions quickly through clearer messaging, better hierarchy, cleaner navigation, and stronger trust signals."
	}

	if contains(userLower, []string{"built to persuade not decorate", "persuade not decorate", "what persuade mean"}) {
		return "It means the site is designed to help the right visitor understand the offer, trust the brand, and take action. The goal is not visual noise or trend-chasing. It is clarity, authority, and conversion."
	}

	// CUSTOM WEBSITE BUILD - POSITIONING AND MESSAGING
	if contains(userLower, []string{"positioning message clarity", "what is message clarity", "positioning clarity definition"}) {
		return "It means making your offer easy to understand. A good site should quickly answer: who you help, what problem you solve, why you are different, and what the next step is. Because users scan pages, concise structure and descriptive headings improve comprehension and navigation."
	}

	if contains(userLower, []string{"why positioning important web design", "positioning importance", "messaging importance design"}) {
		return "Because even a visually polished site underperforms if the visitor cannot tell what the business does or why it matters. Messaging clarity usually has more impact than decorative design."
	}

	if contains(userLower, []string{"help shape copy", "copy writing", "messaging help"}) {
		return "Yes. Based on your offer, that is a core part of the service. The build includes service framing, page hierarchy, and messaging structure so the site explains the business clearly."
	}

	// CUSTOM WEBSITE BUILD - CREDIBILITY AND CONVERSION
	if contains(userLower, []string{"website build credibility", "how build credibility", "credibility through design"}) {
		return "A site builds credibility through design quality, clarity, current and comprehensive content, and transparent structure. Design quality and up-front clarity are key credibility factors, and first impressions strongly shape how relevant and trustworthy a site feels."
	}

	if contains(userLower, []string{"what is conversion path", "conversion path definition", "how conversion works"}) {
		return "A conversion path is the route a visitor takes from interest to action. On a professional services site, that usually means: understand the offer, see proof or trust signals, find the right service, take the next step through inquiry."
	}

	if contains(userLower, []string{"conversion first architecture", "what is conversion architecture", "conversion focused"}) {
		return "It means the site's navigation, page order, calls to action, and section hierarchy are designed to guide qualified visitors toward action instead of making them hunt for answers."
	}

	if contains(userLower, []string{"aggressive marketing tactics", "hard sell", "do i need aggressive"}) {
		return "No. For premium services, credibility, clarity, and reduced friction are often more effective than aggressive tactics. The emphasis is calm, high-trust conversion rather than flashy selling."
	}

	// CUSTOM WEBSITE BUILD - DESIGN AND UX
	if contains(userLower, []string{"premium visual system", "what is visual system", "visual system definition"}) {
		return "A premium visual system is the consistent use of typography, spacing, hierarchy, color, components, and imagery so the site feels coherent and high quality across pages."
	}

	if contains(userLower, []string{"typography spacing matter", "why typography important", "spacing importance"}) {
		return "Because they affect readability, scanning, and perceived professionalism. Users judge trust and usability partly from how clean and organized a site feels on first impression."
	}

	if contains(userLower, []string{"what is hierarchy", "hierarchy definition", "page hierarchy"}) {
		return "Hierarchy is how the page shows what matters most first. Clear headings, spacing, contrast, and layout help visitors understand the page quickly."
	}

	if contains(userLower, []string{"do users really scan", "user behavior scanning", "scan pages"}) {
		return "Yes. Long-standing usability research found that users rarely read web pages word for word and instead scan for relevant words, headings, and cues."
	}

	// CUSTOM WEBSITE BUILD - PERFORMANCE AND TECHNICAL
	if contains(userLower, []string{"performance minded build", "what does performance mean", "performance focused"}) {
		return "It means the site is built to load quickly, stay responsive, and avoid visual instability or bloated code. Load performance, responsiveness, and visual stability are central quality signals for user experience."
	}

	if contains(userLower, []string{"why site speed matters", "speed importance", "fast loading"}) {
		return "Because faster, more responsive sites keep users engaged better than slow ones. Websites that load quickly and respond promptly retain users better than sites that feel sluggish."
	}

	if contains(userLower, []string{"what are core web vitals", "core web vitals", "web vitals definition"}) {
		return "Core Web Vitals are a set of user-experience metrics that evaluate how fast a page feels, how responsive it is, and how visually stable it remains while loading."
	}

	if contains(userLower, []string{"work well on mobile", "mobile responsive", "mobile performance"}) {
		return "That should be a baseline requirement. Sites should be fast, accessible, secure, and work on all devices."
	}

	if contains(userLower, []string{"semantic html", "what is semantic", "semantic structure"}) {
		return "Semantic HTML means using the right structural elements—headings, navigation, main regions, forms, lists, and sections—so content is easier for browsers, assistive technology, and search engines to understand."
	}

	if contains(userLower, []string{"why accessibility matters", "accessibility importance", "accessible design"}) {
		return "Accessibility improves navigation and usability for more users and helps create a structurally stronger site. Headings, semantic regions, and landmarks are key parts of accessible navigation."
	}

	// CUSTOM WEBSITE BUILD - SEO
	if contains(userLower, []string{"seo ready foundation", "what seo ready mean", "seo foundation definition"}) {
		return "It means the site is built so search engines can more easily crawl, understand, and index the content. Descriptive titles, useful headings, good internal linking, crawlable structure, and technically sound pages all support search visibility."
	}

	if contains(userLower, []string{"does design help seo", "design and seo", "seo benefits design"}) {
		return "Indirectly, yes. Better structure, clearer headings, better internal linking, stronger page titles, faster loading, and mobile-friendly implementation all support search visibility and user experience."
	}

	if contains(userLower, []string{"headings titles seo", "do headings matter seo", "titles seo"}) {
		return "Yes. Clear, descriptive title links and site structure matter for search, and descriptive headings help both readers and search engines understand content."
	}

	if contains(userLower, []string{"need plugins seo", "seo plugins", "plugins for seo"}) {
		return "Not necessarily. A clean codebase, good metadata, sound heading structure, strong content hierarchy, and sensible internal links can provide an SEO-ready foundation without bloated plugins."
	}

	// CUSTOM WEBSITE BUILD - PROCESS QUESTIONS
	if contains(userLower, []string{"engagement actually include", "what engagement includes", "what does engagement cover"}) {
		return "Based on our approach, it includes: positioning and message clarity, custom visual system, conversion-focused site architecture, responsive and performance-minded build, SEO-ready technical structure, and launch and handoff."
	}

	if contains(userLower, []string{"only design or also build", "do you build", "design and development"}) {
		return "This is both strategy and implementation: messaging, structure, design, development, launch, and handoff."
	}

	if contains(userLower, []string{"maintainable after launch", "after launch maintenance", "site maintenance"}) {
		return "Yes. The site is organized well enough that your team can update it, extend it, or maintain it without creating chaos. We provide clean file structure, reusable components, and guidance for updates without breaking the design."
	}

	if contains(userLower, []string{"what is handoff", "handoff definition", "handoff meaning"}) {
		return "Handoff means the final site is organized well enough that your team can update it, extend it, or maintain it without chaos."
	}

	// CUSTOM WEBSITE BUILD - OBJECTIONS AND COMPARISONS
	if contains(userLower, []string{"why not template", "template question", "template not enough"}) {
		return "Templates can be useful for simpler needs, but they often come with generic structure and messaging constraints. A custom build is better when credibility, positioning, and conversion quality matter."
	}

	if contains(userLower, []string{"custom overkill", "is custom necessary", "overkill for smaller firm"}) {
		return "Not always. If the site is a major trust signal and lead-generation asset, a stronger custom build can matter even more for a smaller or premium firm because the website often acts as the first serious impression."
	}

	if contains(userLower, []string{"just look good actually help convert", "conversion results", "will this convert"}) {
		return "The goal is both. The service is positioned around trust, messaging clarity, navigation, and inquiry flow—not surface-level aesthetics."
	}

	if contains(userLower, []string{"hire freelance designer", "freelance vs agency", "freelance designer"}) {
		return "That depends on the need. A freelance designer may deliver visuals, but this service is broader: positioning, service narrative, architecture, conversion logic, responsive implementation, SEO structure, and handoff."
	}

	// CUSTOM WEBSITE BUILD - TECHNICAL DETAIL QUESTIONS
	if contains(userLower, []string{"build accessibility start", "accessibility from start", "accessible from beginning"}) {
		return "Yes. That's the right approach. Using headings, semantic structure, and landmarks early rather than patching accessibility in later creates a stronger foundation."
	}

	if contains(userLower, []string{"browser assistive tech", "screen reader", "assistive technology"}) {
		return "Yes. Semantic HTML and landmarks improve programmatic structure and navigation for assistive technologies."
	}

	if contains(userLower, []string{"avoid bloated builds", "bloated code", "unnecessary features"}) {
		return "By using clear structure, reusable components, restrained interactions, and a performance-minded implementation rather than piling on unnecessary plugins or effects."
	}

	// TEMPLATE-BASED BUILD - BEGINNER QUESTIONS
	if contains(userLower, []string{"what is template based build", "template based definition", "what is a template build"}) {
		return "A template-based build uses a proven site foundation, then customizes the structure, messaging, and visual polish so the final site feels credible and tailored rather than generic. It is a faster path than a fully bespoke build, but still focused on clarity, trust, and conversion."
	}

	if contains(userLower, []string{"different from cheap template", "cheap template vs", "cheap template site"}) {
		return "A cheap template site usually stops at installing a theme and swapping in content. A stronger template-based build reshapes the hierarchy, rewrites the messaging, improves typography and spacing, and tunes the flow so visitors understand the offer and know what to do next. Since users tend to scan pages, that structural work matters a lot."
	}

	if contains(userLower, []string{"who is template best for", "template for who", "best suited for template"}) {
		return "It is best for firms that need to launch faster, want a polished high-trust site, do not need a fully custom system yet, and still care about positioning, messaging, and conversion."
	}

	if contains(userLower, []string{"why choose template over custom", "template instead of custom", "choose template build"}) {
		return "Usually because they want a strong site sooner and at a lower cost, while still avoiding the look and feel of a generic quick-build website."
	}

	// TEMPLATE-BASED BUILD - STRATEGY AND MESSAGING
	if contains(userLower, []string{"speed with standards", "what does speed with standards", "speed and standards"}) {
		return "It means the site launches faster because it starts from a strong foundation, but it still follows good standards for messaging clarity, design consistency, mobile behavior, accessibility structure, and technical basics."
	}

	if contains(userLower, []string{"help with messaging", "messaging help template", "does messaging help"}) {
		return "Yes. Messaging is a core part of the engagement. That includes clarifying what the company does, shaping the service narrative, and placing calls to action where visitors are most likely to decide."
	}

	if contains(userLower, []string{"why messaging big deal", "messaging importance website", "why messaging matters"}) {
		return "Because visitors scan quickly. If they cannot understand the offer in a few seconds, trust drops and conversion suffers. People scan web pages and concise, scannable writing improves usability."
	}

	if contains(userLower, []string{"make instantly understandable", "instantly understandable", "quick understanding"}) {
		return "It means the site should quickly answer: who you help, what you do, why it matters, and what the next step is."
	}

	// TEMPLATE-BASED BUILD - CREDIBILITY AND CONVERSION
	if contains(userLower, []string{"template feel premium", "premium template", "can template be premium"}) {
		return "Yes. A site feels premium when the structure is disciplined, the messaging is clear, the typography and spacing are coherent, and the experience feels intentional. First impressions strongly affect credibility and perceived usability, so polish and hierarchy matter even on a template foundation."
	}

	if contains(userLower, []string{"what is conversion focused", "conversion focused site", "conversion focused definition"}) {
		return "It is a site built to guide the right visitor toward action, usually through: clear navigation, strong page hierarchy, well-placed calls to action, reduced confusion, and a simple inquiry path."
	}

	if contains(userLower, []string{"help reduce what do you confusion", "what do you do confusion", "reduce confusion"}) {
		return "Yes. That is one of the biggest advantages of this service. The site structure and copy are deliberately shaped to make the offer clearer and easier to trust."
	}

	if contains(userLower, []string{"why cta matter", "why calls to action", "cta importance"}) {
		return "Because even interested users often need direction. A good CTA appears where a user is ready to act, rather than forcing them to hunt for the next step."
	}

	// TEMPLATE-BASED BUILD - DESIGN AND UX
	if contains(userLower, []string{"visual system brand polish", "what is visual system", "brand polish definition"}) {
		return "It means the typography, spacing, colors, components, imagery, and layout all feel consistent and intentional across the site."
	}

	if contains(userLower, []string{"spacing typography matter", "why spacing matter", "typography importance template"}) {
		return "They affect readability, scanning, and perceived quality. Design quality influences how credible and usable a site feels on first impression."
	}

	if contains(userLower, []string{"template designed well", "can template be well designed", "design quality template"}) {
		return "Yes. A strong template can be an excellent starting point if the layout is good and the final site is refined carefully instead of treated like plug-and-play software."
	}

	// TEMPLATE-BASED BUILD - TECHNICAL QUESTIONS
	if contains(userLower, []string{"responsive performance build", "what responsive performance", "responsive and performance"}) {
		return "It means the site is built to work well across devices, load quickly, remain stable while loading, and feel responsive to users. Core Web Vitals focus on perceived load speed, responsiveness, and visual stability."
	}

	if contains(userLower, []string{"why performance matters", "performance importance", "performance matters"}) {
		return "Because poor performance hurts user experience, especially on mobile devices with weaker CPU and network conditions. Users tolerate poor performance only for so long before abandoning a site."
	}

	if contains(userLower, []string{"semantic html template", "what semantic html", "semantic html definition"}) {
		return "Semantic HTML means using page structure and elements in ways that clearly communicate organization, such as headings, navigation, main content regions, and forms. Landmarks and structured page regions help users of assistive technology understand and navigate a page more easily."
	}

	if contains(userLower, []string{"accessibility part", "why accessibility", "accessibility importance template"}) {
		return "Because accessible structure improves navigation, readability, and maintainability for everyone, not just users of assistive tech. Good page structure also helps search engines better index content."
	}

	// TEMPLATE-BASED BUILD - SEO QUESTIONS
	if contains(userLower, []string{"basic seo essentials", "what seo essentials", "seo essentials definition"}) {
		return "That means the site includes core on-page fundamentals like: descriptive page titles, meta descriptions, logical heading structure, crawlable internal links, and useful text around important content."
	}

	if contains(userLower, []string{"why mobile seo", "mobile important seo", "mobile seo importance"}) {
		return "Because Google uses mobile-first indexing, meaning the mobile version of a site is the primary version used for indexing and ranking."
	}

	if contains(userLower, []string{"template based rank", "can template rank", "template seo ranking"}) {
		return "Yes, if the content is helpful, the site structure is clear, the pages are crawlable, the titles and headings are descriptive, and the mobile and performance experience is solid."
	}

	// TEMPLATE-BASED BUILD - PROCESS QUESTIONS
	if contains(userLower, []string{"pages usually included", "what pages included", "included pages template"}) {
		return "This usually includes: home, services overview, primary service detail pages, and contact or inquiry flow."
	}

	if contains(userLower, []string{"launch handoff mean", "what is launch handoff", "launch handoff definition"}) {
		return "It means the site is reviewed before launch for issues like spacing, broken links, and readability, then handed off with simple notes so future updates stay clean."
	}

	if contains(userLower, []string{"able update site later", "update after launch", "site maintenance template"}) {
		return "Yes. The goal is a clean, maintainable foundation you can build on later without starting over."
	}

	if contains(userLower, []string{"scale into custom", "scale to custom", "can template scale"}) {
		return "Yes. This is positioned as a foundation that can grow later into something more custom without rebuilding from scratch."
	}

	// TEMPLATE-BASED BUILD - COMPARISON AND OBJECTIONS
	if contains(userLower, []string{"different from custom build", "template vs custom build", "custom vs template build"}) {
		return "A custom build starts from scratch around your exact business, while a template-based build starts from a strong base and customizes the parts that matter most. The template option is usually faster and less expensive."
	}

	if contains(userLower, []string{"still feel unique", "feel unique template", "template unique"}) {
		return "It should feel tailored in the places that matter most: messaging, hierarchy, typography, brand polish, and conversion flow."
	}

	if contains(userLower, []string{"just shortcut", "is shortcut template", "shortcut approach"}) {
		return "It is a faster path, but not a careless one. The point is to launch efficiently without sacrificing trust, clarity, or quality."
	}

	if contains(userLower, []string{"why not buy template ourselves", "do it ourselves template", "self run template"}) {
		return "Because most self-run template sites stop at surface customization. The bigger value is in the strategy: structure, messaging, hierarchy, polish, and conversion flow."
	}

	// WEBSITE REDESIGN - BEGINNER QUESTIONS
	if contains(userLower, []string{"what is website redesign", "website redesign definition", "what is redesign"}) {
		return "A website redesign improves an existing site's structure, messaging, visual quality, and user flow without unnecessarily rebuilding everything from scratch. The goal is to keep what already works, remove friction, and make the site clearer and more credible."
	}

	if contains(userLower, []string{"redesign different from custom build", "redesign vs custom build", "redesign vs full custom"}) {
		return "A redesign starts with an existing site and improves it. A custom build starts from zero. A redesign is often right when the business already has useful content, pages, or SEO value, but the site feels dated, unclear, or harder to trust."
	}

	if contains(userLower, []string{"why not just look newer", "why not make it look newer", "why not freshen it up"}) {
		return "Because a redesign that only changes visuals can make the site worse. Good redesign work improves clarity, structure, and conversion—not just appearance."
	}

	if contains(userLower, []string{"what problem redesign solve", "problem redesign solves", "why would i need redesign"}) {
		return "Usually one or more of these: the site feels dated, the messaging is unclear, navigation feels too heavy, mobile experience is weak, inquiries are lower than they should be, or the business has evolved but the site has not."
	}

	// WEBSITE REDESIGN - STRATEGY AND MESSAGING
	if contains(userLower, []string{"refine don't reinvent", "refine not reinvent", "what refine mean"}) {
		return "It means preserving the pieces that already help trust or conversion, then improving the structure and presentation around them. The point is disciplined improvement, not change for its own sake."
	}

	if contains(userLower, []string{"know what stay change", "what should change", "what should stay"}) {
		return "You start with a content and structure audit. That means identifying: which pages already perform well, which content builds trust, where users likely get confused, and where navigation or page hierarchy is slowing decisions."
	}

	if contains(userLower, []string{"what is content structure audit", "content audit definition", "structure audit"}) {
		return "It is a review of how the site is organized, what each page is trying to do, and where users may be dropping off or losing clarity. A good redesign usually starts there rather than jumping straight into visuals."
	}

	if contains(userLower, []string{"why messaging refresh", "messaging refresh importance", "messaging refresh redesign"}) {
		return "Because visitors scan, not read line by line. Users scan web pages and benefit from concise, scannable writing and structure. That means stronger messaging often improves usability and conversion as much as visual design does."
	}

	if contains(userLower, []string{"make instantly understandable", "instantly understandable", "quick understanding"}) {
		return "It means the site should quickly answer: who you help, what you do, why it matters, and what the next step is."
	}

	// WEBSITE REDESIGN - CREDIBILITY AND CONVERSION
	if contains(userLower, []string{"redesign improve credibility", "how improve credibility", "credibility through redesign"}) {
		return "By improving first impressions, reducing confusion, strengthening visual consistency, and making the business feel more current and trustworthy. First impressions strongly influence how users perceive relevance, credibility, and even usability."
	}

	if contains(userLower, []string{"what is credibility first redesign", "credibility first", "credibility redesign"}) {
		return "A credibility-first redesign focuses on: clearer hierarchy, stronger messaging, more polished visual detail, cleaner navigation, and lower-friction inquiry flow. It is less about looking trendy and more about making the business easier to trust."
	}

	if contains(userLower, []string{"redesign affect conversion", "conversion improvement redesign", "how affect conversion"}) {
		return "It can improve conversion by making it easier for visitors to understand the offer, find the right page, and take action. Better CTA placement, cleaner navigation, and more intentional page flow usually reduce friction."
	}

	if contains(userLower, []string{"what is conversion ready inquiry", "conversion inquiry flow", "inquiry flow definition"}) {
		return "It means the path to contact or inquire is straightforward, visible, and low-friction. The user should not have to wonder where to click next or whether they are on the right page."
	}

	// WEBSITE REDESIGN - DESIGN AND UX
	if contains(userLower, []string{"what is visual system upgrade", "visual system upgrade", "visual upgrade"}) {
		return "It means improving typography, spacing, layout rhythm, component consistency, and overall polish across the site so it feels coherent and premium."
	}

	if contains(userLower, []string{"spacing hierarchy matter", "why spacing hierarchy", "spacing hierarchy importance"}) {
		return "Because they control readability and scannability. Headings communicate page organization, and browsers and assistive technologies use them for in-page navigation. Users scan for structure and cues."
	}

	if contains(userLower, []string{"cleaner reading rhythm", "what is reading rhythm", "reading rhythm definition"}) {
		return "It means the page feels easier to move through visually. Sections are easier to distinguish, headings are more meaningful, spacing is more intentional, and the content feels less dense or chaotic."
	}

	if contains(userLower, []string{"redesign improve trust without flashy", "trust without flashy", "trust without trendy"}) {
		return "Yes. In many premium-service environments, restrained, disciplined design communicates more trust than louder or trendier design."
	}

	// WEBSITE REDESIGN - TECHNICAL QUESTIONS
	if contains(userLower, []string{"performance mobile polish", "what performance mobile", "performance and mobile"}) {
		return "It means improving how the site behaves on real devices: load speed, responsiveness, stability, image handling, font loading, and layout behavior. Core Web Vitals focus on perceived load speed, responsiveness, and visual stability."
	}

	if contains(userLower, []string{"why performance matters redesign", "performance importance", "performance matters"}) {
		return "Because users experience the redesign through performance as much as visuals. A site that looks cleaner but loads poorly or shifts around on mobile will still feel lower quality."
	}

	if contains(userLower, []string{"responsive tuning", "what responsive tuning", "tuning definition"}) {
		return "It means the redesigned site is tested and refined for how layouts, text, navigation, and interactions behave across screen sizes—especially on mobile, where constraints are tighter."
	}

	if contains(userLower, []string{"why semantic structure matter", "semantic structure importance", "semantic structure"}) {
		return "Semantic structure and landmarks help users of assistive technologies understand and navigate the page. Headings and structural markup also improve content organization more broadly."
	}

	// WEBSITE REDESIGN - SEO QUESTIONS
	if contains(userLower, []string{"can redesign hurt seo", "redesign hurt search", "seo risk redesign"}) {
		return "Yes, if it changes structure carelessly, removes valuable content, breaks links, or weakens headings and metadata. That is why a disciplined redesign matters."
	}

	if contains(userLower, []string{"redesign without losing seo", "preserve seo", "keep seo value"}) {
		return "By preserving useful content, keeping strong pages in mind, maintaining logical site structure, and carrying forward important SEO elements like descriptive titles, headings, and crawlable links."
	}

	if contains(userLower, []string{"basic seo essentials redesign", "seo essentials", "essential seo"}) {
		return "Usually: descriptive page titles, strong heading hierarchy, logical internal links, readable page structure, and preserved or improved crawlability."
	}

	if contains(userLower, []string{"mobile matter search", "mobile seo", "mobile indexing"}) {
		return "Yes. Google uses mobile-first indexing, which means the mobile version is the primary version considered for indexing and ranking."
	}

	// WEBSITE REDESIGN - PROCESS QUESTIONS
	if contains(userLower, []string{"redesign engagement include", "what engagement includes", "what does redesign include"}) {
		return "Based on our approach, it includes: content and structure audit, messaging refresh, visual system upgrade, navigation and conversion flow refinement, performance and mobile polish, and launch QA and handoff."
	}

	if contains(userLower, []string{"rebuild everything", "rebuild from scratch", "rebuild all"}) {
		return "Not necessarily. The point is to improve what needs improvement while keeping what already supports trust, clarity, or discoverability."
	}

	if contains(userLower, []string{"keep existing content", "preserve content", "keep content"}) {
		return "Some of it, yes—especially if it already performs well or builds credibility. But key sections may be rewritten for clarity and authority."
	}

	if contains(userLower, []string{"what is launch qa", "launch qa definition", "qa before launch"}) {
		return "It is the final review before launch to catch issues in spacing, broken links, readability, responsiveness, and other polish details."
	}

	if contains(userLower, []string{"what is handoff redesign", "handoff definition", "handoff meaning"}) {
		return "Handoff means you receive a clean, consistent final site plus guidance so future edits stay aligned instead of slowly degrading the design."
	}

	// WEBSITE REDESIGN - OBJECTION AND COMPARISON QUESTIONS
	if contains(userLower, []string{"why redesign instead of scratch", "redesign vs scratch", "rebuilding from scratch"}) {
		return "Because sometimes the business already has useful structure, valuable content, or search equity worth preserving. A redesign can modernize the experience with less disruption."
	}

	if contains(userLower, []string{"know if redesign or custom", "redesign or custom build", "when redesign when custom"}) {
		return "A redesign is usually right when the site has something worth preserving but needs stronger clarity, polish, or conversion flow. A custom build is better when the current site is too limited structurally or no longer reflects the business at all."
	}

	if contains(userLower, []string{"redesign disrupt what works", "will disrupt", "disruption risk"}) {
		return "It should not if done carefully. The goal is explicitly to avoid unnecessary disruption."
	}

	if contains(userLower, []string{"redesign improve trust without changing", "trust without change", "improve without changing"}) {
		return "Yes. Often the biggest improvements come from hierarchy, copy clarity, navigation, and polish rather than total reinvention."
	}

	if contains(userLower, []string{"case study", "success story", "client", "results", "example", "proven", "roi"}) {
		return "We've helped numerous companies transform their data strategy and achieve measurable results. From increasing revenue visibility to improving decision-making speed, our clients see real business impact. For specific examples relevant to your industry, please connect with us at inquire page."
	}

	// UX / UI DESIGN - BEGINNER QUESTIONS
	if contains(userLower, []string{"what is ux ui design", "ux ui definition", "what is ux and ui"}) {
		return "UX design focuses on how a product works: the flow, structure, clarity, and ease of use. UI design focuses on how the product looks and communicates visually: layout, typography, spacing, components, and polish. Together, they make a product easier to understand, easier to trust, and easier to use."
	}

	if contains(userLower, []string{"why ux ui matter", "why does design matter", "ux ui importance"}) {
		return "Because users make judgments quickly. If a product feels confusing, cluttered, or inconsistent, trust drops and task completion suffers. Good UX/UI reduces friction and makes important actions feel obvious. Users scan rather than read word-for-word, which makes hierarchy and clarity especially important."
	}

	if contains(userLower, []string{"what problem ux ui solve", "problem ux ui solves", "why would i need ux ui"}) {
		return "Usually one or more of these: users hesitate, flows feel harder than they should, screens look cluttered, conversion is weaker than expected, the product works but does not feel intuitive, or the interface doesn't reflect the quality of the business."
	}

	if contains(userLower, []string{"make complex products simple", "complex products simple", "simplify complexity"}) {
		return "It means reducing mental load. The product may still be sophisticated underneath, but the interface should help users know where to look, what matters, and what to do next without overexplaining."
	}

	// UX / UI DESIGN - UX QUESTIONS
	if contains(userLower, []string{"what is ux simple terms", "ux definition", "what is user experience"}) {
		return "UX is the structure of the experience. It covers: user journeys, task flows, navigation, friction points, screen sequence, and how easily someone can complete an important task."
	}

	if contains(userLower, []string{"what is flow", "flow definition", "what does flow mean"}) {
		return "A flow is the path a user takes through the product to complete something important, like signing up, submitting information, reviewing data, or making a decision."
	}

	if contains(userLower, []string{"what is friction", "friction definition", "friction in design"}) {
		return "Friction is anything that makes a task harder than it needs to be: too many steps, confusing labels, weak hierarchy, uncertainty about the next action, or inconsistent behavior across screens."
	}

	if contains(userLower, []string{"how reduce friction", "reduce friction", "friction reduction"}) {
		return "Usually by: simplifying paths, clarifying labels and actions, reducing unnecessary choices, making key actions more visible, and improving layout and visual hierarchy."
	}

	if contains(userLower, []string{"what are user journeys", "user journey definition", "user journey mapping"}) {
		return "User journeys map what a user is trying to do from beginning to end. They help identify where confusion, hesitation, or drop-off is happening."
	}

	// UX / UI DESIGN - UI QUESTIONS
	if contains(userLower, []string{"what is ui simple terms", "ui definition", "what is user interface"}) {
		return "UI is the visual system of the product: typography, spacing, grids, buttons, inputs, cards, colors, and other interface patterns."
	}

	if contains(userLower, []string{"why typography spacing matter", "typography importance", "spacing importance"}) {
		return "Because they affect readability, scannability, and perceived quality. Clear headings and organized structure help users navigate content. Headings communicate page organization for browsers and assistive technologies."
	}

	if contains(userLower, []string{"what is hierarchy design", "hierarchy definition", "visual hierarchy"}) {
		return "Hierarchy is how the design shows what matters most first. A good hierarchy makes the important content and actions stand out immediately, which matters because users tend to scan pages for cues rather than read everything closely."
	}

	if contains(userLower, []string{"what is ui system", "ui system definition", "design system"}) {
		return "A UI system is a consistent set of rules for layout, spacing, typography, components, and interaction patterns so screens feel unified instead of pieced together."
	}

	// UX / UI DESIGN - PRODUCT AND BUSINESS QUESTIONS
	if contains(userLower, []string{"ux ui help conversion", "conversion improvement design", "how affect conversion"}) {
		return "It improves conversion by making the desired action clearer and easier. Better hierarchy, clearer CTAs, cleaner flows, and less hesitation usually lead to stronger follow-through."
	}

	if contains(userLower, []string{"ux ui help trust", "trust through design", "design and trust"}) {
		return "Design communicates competence. A calm, consistent interface feels more reliable than one that feels messy or improvised. The interface should be obvious to use and worth trusting."
	}

	if contains(userLower, []string{"ux ui help task completion", "task completion speed", "faster completion"}) {
		return "Yes. Cleaner flows and clearer decision points usually reduce time spent figuring out what to do. That is one reason thoughtful hierarchy and friction reduction matter so much."
	}

	if contains(userLower, []string{"only websites", "website only", "apply to products"}) {
		return "No. This applies to products, dashboards, apps, internal tools, portals, and service experiences where users need to navigate information and take action."
	}

	// UX / UI DESIGN - TECHNICAL AND ADVANCED QUESTIONS
	if contains(userLower, []string{"developer ready specs", "what developer ready", "specs meaning"}) {
		return "It means the design is delivered in a way that developers can build from: clear layouts, reusable components, states, annotations, spacing rules, and implementation notes."
	}

	if contains(userLower, []string{"consistency across screens", "why consistency matter", "consistency importance"}) {
		return "Because consistency reduces cognitive load. When buttons, layouts, and component behavior change unpredictably, users have to re-learn the interface on every screen."
	}

	if contains(userLower, []string{"semantic structure", "what semantic structure", "semantic meaning"}) {
		return "Semantic structure means using meaningful page organization and structural elements so content is easier for browsers, search engines, and assistive technologies to interpret. Landmarks and structured sections help users navigate and understand pages more efficiently."
	}

	if contains(userLower, []string{"mobile behavior ux", "mobile important", "mobile first"}) {
		return "Mobile should be considered early, not patched in later. Google uses mobile-first indexing for search, and developer guidance recommends sites be fast and work on all devices."
	}

	if contains(userLower, []string{"performance relate ux", "performance and design", "performance impact"}) {
		return "Performance is part of the experience. Core Web Vitals measure things like responsiveness and visual stability. Unexpected layout shifts, for example, directly hurt usability."
	}

	if contains(userLower, []string{"what is visual stability", "visual stability definition", "layout shift"}) {
		return "Visual stability means the interface does not jump around unexpectedly while loading or updating. This is an important quality signal for user experience."
	}

	// UX / UI DESIGN - PROCESS QUESTIONS
	if contains(userLower, []string{"ux ui engagement include", "what engagement includes", "what does engagement cover"}) {
		return "Based on our approach, it includes: user journeys and key path mapping, friction review and flow simplification, UI system work, high-end visual polish, conversion tuning, Figma designs and when useful a clickable prototype, developer-ready handoff, and iteration and polish pass."
	}

	if contains(userLower, []string{"clients actually receive", "what will i get", "deliverables"}) {
		return "Typically: designed or redesigned screens, a clearer flow structure, component and layout consistency, a prototype when needed, specs and notes for development, and guidance for future updates."
	}

	if contains(userLower, []string{"only redesign existing", "design new screens", "new designs"}) {
		return "Both. The service covers screen designs or redesigns in Figma."
	}

	if contains(userLower, []string{"what is clickable prototype", "prototype definition", "interactive mockup"}) {
		return "A clickable prototype is an interactive mockup that lets stakeholders review the experience before development."
	}

	// UX / UI DESIGN - OBJECTION AND COMPARISON QUESTIONS
	if contains(userLower, []string{"why not prettier", "just make prettier", "visual polish only"}) {
		return "Because visual polish without UX clarity often fails. A better-looking interface can still confuse users if the structure and hierarchy are weak."
	}

	if contains(userLower, []string{"different from branding", "branding vs ux", "branding differences"}) {
		return "Branding shapes overall identity. UX/UI shapes how users move through the actual product and how clearly the interface communicates."
	}

	if contains(userLower, []string{"need if already works", "product works fine", "already functional"}) {
		return "Possibly, yes. A product can function technically while still creating hesitation, inefficiency, or trust issues for users."
	}

	if contains(userLower, []string{"know if issue ux ui", "how identify problem", "identify ux issue"}) {
		return "Common signs: users need too much explanation, important actions are being missed, flows feel heavier than they should, internal teams disagree about what users should do next, or the product feels less premium than the business behind it."
	}

	// OFF-TOPIC HANDLER - catch before team/experience handler
	offTopicPatterns := []string{"what did you eat", "eat today", "lunch", "dinner", "breakfast", "what's for dinner", "pizza", "burger", "weather today", "weather forecast", "movie", "funny", "joke", "favorite pet", "cat", "dog", "vacation", "travel", "recipe", "cooking", "what do you like"}
	if contains(userLower, offTopicPatterns) && !contains(userLower, []string{"marault", "data", "service", "business", "analytics", "audit", "dashboard", "build", "design", "website"}) {
		return "That's a fun question, but I'm specifically here to help with Marault Intelligence services. Feel free to ask me about our services, team, approach, or anything related to data analytics and business technology."
	}

	if contains(userLower, []string{"who are you", "background", "qualification", "caroline", "lindsey", "consultant", "expertise about you", "your expertise"}) {
		return "Our team includes senior consultants with deep expertise in data strategy, analytics, and business technology. Meet the team at executive team page. We bring elite training and expertise from leading companies and understand the complexities of data transformation. How can we help you?"
	}

	// SUPPORT & MAINTENANCE
	if contains(userLower, []string{"support", "maintenance", "ongoing", "after launch", "training", "help", "update"}) {
		return "We offer ongoing support and maintenance for all our deliverables. This includes training your team, handling updates, troubleshooting, and optimizations. Support packages can be customized based on your needs. Let's discuss your support requirements at inquire page."
	}

	// METHODOLOGY & APPROACH
	if contains(userLower, []string{"approach", "philosophy", "methodology", "process", "how", "do you work"}) {
		return "We believe in making data accessible and actionable for all stakeholders. Our approach involves deep discovery, strategic planning, implementation, and ongoing optimization. Learn more about our philosophy and approach pages."
	}

	// GETTING STARTED
	if contains(userLower, []string{"start", "getting started", "next step", "begin", "ready", "let's work", "ready to start", "let's talk", "move forward", "take action", "want help"}) {
		return "Great! Ready to transform your business? At Marault Intelligence, we work through a simple process: 1) Schedule a consultation call (15-30 minutes), 2) We'll learn about your specific challenges and goals, 3) We'll propose a tailored solution with timeline and scope, 4) You decide and we execute with regular check-ins. Every project is unique in scope and timeline. Start the conversation at inquire page—our team responds within 24 hours."
	}

	// CONTACT & COMMUNICATION
	if contains(userLower, []string{"contact", "call", "email", "phone", "reach", "schedule", "meeting", "call us", "how to reach", "get in touch", "connect with"}) {
		return "You can reach Marault Intelligence in several ways: Fill out our inquiry form at inquire page for service inquiries or to schedule a consultation. Visit contact page for general questions. Our team responds within 24 hours. You can also reach our consultants Caroline and Lindsey directly through our contact page. We're available to discuss your challenges and explore solutions."
	}

	// ALL SERVICES OVERVIEW
	if contains(userLower, []string{"services", "offer", "what do you", "what services", "capabilities", "what can you help with", "company services", "what you offer", "service offerings"}) {
		return "Marault Intelligence specializes in transforming data and design into competitive advantages. We offer 8 core services across two areas: Business Intelligence (Data Visibility Audits, Revenue & Customer Analytics, Executive Dashboards & Reporting, Forecasting & Decision Modeling, Private Client Analytics) and Web & Experience (Custom Website Builds, Template-Based Builds, Website Redesigns, UX/UI Design). Whether your challenge is data clarity, decision speed, wealth management, or digital presence, we have a solution. Visit services page to explore each one in detail."
	}

	// HELLO & GENERAL HELP
	if contains(userLower, []string{"hello", "hi", "hey", "help", "support", "question", "info", "what's this", "who are you", "tell me about", "about marault", "about the company"}) {
		return "Hi there! I'm the Marault Intelligence chatbot, here to help. Marault Intelligence transforms data confusion into clarity and designs conversion-focused experiences for strategy-driven companies. We work with firms across Finance, SaaS, E-Commerce, Professional Services, and more. You can ask me about our services, team, pricing, timeline, or anything else. Or explore at services page, approach page, philosophy page, or executive team page. What would you like to know?"
	}

	// ==========================================
	// CONTEXT-AWARE RESPONSES (Using Memory)
	// ==========================================

	// If asking about a specific service after previous context
	if contains(userLower, []string{"tell me more", "more details", "more info", "expand on that", "elaborate", "go deeper", "dig deeper", "specific details", "walk me through"}) ||
		contains(userLower, []string{"how does this work", "how would this help us", "can you explain more", "what exactly", "give me more", "more about this"}) {

		// Check if they've been asking about data visibility audits
		if contains(conversationContext, []string{"data visibility audit", "data audit", "audit my data", "understand data"}) {
			return "Great question! A Data Visibility Audit delivers: 1) KPI definitions that your entire organization agrees on, 2) A clean reporting baseline so everyone starts from truth, 3) Data inventory and lineage showing where everything comes from, 4) A data quality scorecard identifying gaps, 5) An executive dashboard with decision-ready metrics, and 6) A stabilization roadmap for ongoing improvement. The goal is to establish one trusted source of truth for decisions. Every project is unique in scope and timeline. At Marault Intelligence, we combine strategy and technology to make this foundation maintainable long-term. Ready to explore this further? Visit services page on data visibility audit."
		}

		// Check if they've been asking about executive dashboards
		if contains(conversationContext, []string{"executive dashboard", "executive reporting", "dashboards", "reporting", "metrics"}) {
			return "Absolutely! Executive Dashboards typically include 5-12 carefully chosen metrics, narrative explanations of what changed and why it matters, drill-down capabilities for deeper exploration, and ongoing governance structures. The key difference from other dashboards is we prioritize clarity and decision-readiness over comprehensive data display. Most teams see a 50% reduction in review time and faster decision-making. At Marault Intelligence, we design these specifically to support executive leadership. Want to learn more? Visit services page on executive dashboards reporting."
		}

		// Check if they've been asking about private client analytics
		if contains(conversationContext, []string{"private client", "wealth", "portfolio", "clients"}) {
			return "Great! For Private Client Analytics, we consolidate accounts, holdings, and risk exposures across multiple institutions into one clear view. That includes performance attribution, risk exposure mapping, scenario analysis, and tax-aware decision support. The engagement typically includes content and structure audit, visual system upgrade, conversion flow refinement, and launch handoff. At Marault Intelligence, we specialize in helping wealth advisors and family offices provide clearer, more confident recommendations to their clients. Ready to explore? Visit services page on private client analytics."
		}

		// Check if they've been asking about forecasting
		if contains(conversationContext, []string{"forecasting", "decision modeling", "predict", "scenario", "planning"}) {
			return "Excellent question! Forecasting & Decision Modeling includes scenario analysis showing how your business could perform under different conditions, sensitivity analysis identifying key drivers, stress testing for risk preparation, and decision frameworks helping leadership choose between options with confidence. The engagement covers discovery, model building, validation, training, and ongoing support. At Marault Intelligence, we help companies move from guesswork to data-backed planning. Most clients see faster decision cycles and better risk preparedness. Ready to dive deeper? Visit services page on forecasting decision modeling."
		}

		// Check if they've been asking about website services
		if contains(conversationContext, []string{"website", "web design", "custom build", "redesign", "template", "ux", "ui"}) {
			return "Perfect! Our web services include strategy (positioning, messaging, hierarchy), design (visual systems, user experience, accessibility), development (clean, performant code), and launch support. Whether it's a custom website build, template-based build, website redesign, or UX/UI design, our approach ensures credibility, clarity, and conversion-ready experiences. At Marault Intelligence, we don't just build websites—we build conversion assets that communicate your brand's authority. Ready to learn more? Visit services page on web design services."
		}
	}

	// If they ask about implementation after discussing a service
	if contains(userLower, []string{"how do we start", "first step", "next steps", "ready to get started", "let's begin", "how to begin", "what's next", "what's the process", "how do we proceed", "get started"}) {
		if contains(conversationContext, []string{"data visibility audit", "audit"}) {
			return "Perfect! The first step is a 15-20 minute discovery call. We'll understand your current reporting challenges, the systems you use, your goals, and your team structure. This helps us scope exactly what you need and create a tailored proposal. At Marault Intelligence, we've worked with companies across industries to fix their data foundations—we know what to look for. Head over to inquire page to schedule your conversation with our team."
		}
		if contains(conversationContext, []string{"executive dashboard", "dashboard"}) {
			return "Great! We start with KPI alignment—a conversation where leadership defines what metrics truly matter for decisions. From there, we design, build, validate, and train your team. At Marault Intelligence, we ensure the dashboard becomes a real tool for decision-making, not just another report. Ready to begin? Visit inquire page and we'll set up your initial alignment session with our consultants."
		}
		if contains(conversationContext, []string{"custom build", "website redesign", "ux ui"}) {
			return "Excellent! We start with a discovery and strategy session to understand your business, positioning, goals, and audience. Then we move into design and development phases with regular check-ins. At Marault Intelligence, we bring both design excellence and conversion strategy to every web project. Head to inquire page to start your project conversation."
		}
	}

	// If they're comparing services
	if contains(userLower, []string{"which should we", "which is better", "difference between", "compared to", "when use", "which one", "both needed"}) {
		if contains(conversationContext, []string{"audit"}) && contains(userLower, []string{"dashboard"}) {
			return "Great instinct! Many companies do both, and here's why: The audit fixes your foundation first (definitions, quality, structure), then the dashboard displays the results beautifully and enforces ongoing governance. Some companies need just the audit if they don't need ongoing dashboards, others need both for complete data strategy. At Marault Intelligence, we help you determine the right path for your specific situation. Let's chat about your needs at inquire page."
		}
		if contains(conversationContext, []string{"website", "web"}) && (contains(userLower, []string{"custom", "template", "redesign"}) || contains(userLower, []string{"which build"})) {
			return "Excellent question! Here's how we typically recommend thinking about it: Choose a Custom Build if you need full control and strategic positioning from scratch. Choose a Template-Based Build if you want a strong site quickly at lower cost. Choose a Redesign if your site has value but needs updating. Choose UX/UI Design if your product feels clunky. At Marault Intelligence, we'll help you pick the right path based on your goals and timeline. Visit inquire page to discuss your project."
		}
	}

	// If they ask about ROI or results
	if contains(userLower, []string{"results", "what happens", "what outcome", "roi", "return", "improve", "benefit", "advantage", "success", "get from this"}) {
		if contains(conversationContext, []string{"data visibility audit", "audit"}) {
			return "Most companies experience: clearer decision-making, faster review cycles (often 50%+ reduction), elimination of conflicting reports, better board and leadership confidence, and a roadmap for ongoing improvement. At Marault Intelligence, we measure success by your team's ability to make faster, more confident decisions backed by reliable data. The timeline varies by scope. Ready to see what's possible for your organization? Visit inquire page."
		}
		if contains(conversationContext, []string{"executive dashboard"}) {
			return "Companies typically see: reduced reporting time by 50%+, faster decision cycles, higher confidence in metrics, alignment across teams, and better strategic focus. At Marault Intelligence, we design dashboards that actually get used and drive action. Your metrics become a real management tool. Visit inquire page to discuss your metrics goals."
		}
		if contains(conversationContext, []string{"website"}) {
			return "Better websites typically result in: improved trust signals in first impressions, clearer understanding of your offer, reduced visitor confusion, and increased inquiry conversion. At Marault Intelligence, we build websites that work as business assets—supporting positioning, reinforcing credibility, and moving qualified prospects toward action. Visit inquire page to start your web project."
		}
	}

	// If they ask "what else" or want to know about related services
	if contains(userLower, []string{"what else", "anything else", "other services", "what about", "other options", "more solutions"}) {
		if contains(conversationContext, []string{"data visibility audit", "audit"}) {
			return "After an audit foundation, many companies expand into Executive Dashboards for ongoing decision support, or Forecasting & Decision Modeling for strategic planning. Some move into Revenue & Customer Analytics for deeper business insights. At Marault Intelligence, we help companies build incrementally—each service strengthens your analytics maturity. What would you like to explore?"
		}
		if contains(conversationContext, []string{"executive dashboard", "dashboard"}) {
			return "Many clients combine dashboards with deeper analytics like Revenue & Customer Analytics or Forecasting & Decision Modeling. Or if you're also thinking about your digital presence, we offer comprehensive web services. At Marault Intelligence, we work best as strategic partners who understand your whole business. What else are you thinking about?"
		}
		if contains(conversationContext, []string{"website", "web design"}) {
			return "If you need help with your data strategy alongside your web presence, we offer comprehensive analytics services. Many companies benefit from both strong digital positioning and data-driven decision-making. At Marault Intelligence, we help companies compete through clarity and credibility. Are you interested in exploring other areas?"
		}
	}

	// If they ask about specific industries or company types
	if contains(userLower, []string{"our industry", "we're a", "we work in", "our company", "company profile", "our business", "finance", "saas", "ecommerce", "healthcare", "professional services"}) {
		return "Marault Intelligence works with strategy-driven companies across Finance, SaaS, E-Commerce, Healthcare, Professional Services, Real Estate, and more. Each industry has unique data and design challenges. At Marault Intelligence, we customize solutions based on your specific business model, complexity, and goals. Tell us more about your industry and challenges at inquire page, and we'll propose solutions tailored to your situation."
	}

	// If they compare us or ask about alternatives
	if contains(userLower, []string{"compared to", "vs other", "why marault", "why choose you", "why not", "differentiate", "unique", "better than", "how different"}) {
		return "At Marault Intelligence, we combine strategic thinking with technical execution. We don't just build dashboards or websites—we build conversion assets and decision-enabling systems. We emphasize clarity, credibility, and measurable outcomes. Our approach involves deep understanding of your business, discipline in execution, and commitment to results that matter. We'd love to discuss how we differ from other firms and why we might be right for you. Visit inquire page to schedule a conversation."
	}

	// CORE MISSION & COMPANY IDENTITY
	// ==========================================

	// CORE MISSION STATEMENT
	if contains(userLower, []string{"mission", "your mission", "what's your mission", "marault mission", "mission statement"}) {
		return "Our core mission: Marault Intelligence helps companies and individuals make better decisions by turning unclear data, unclear positioning, and unclear systems into structured, decision-ready clarity. We believe that clarity is competitive advantage. Most businesses struggle not because they lack information, but because they can't see clearly or act confidently. We fix that."
	}

	// ==========================================
	// BASIC QUESTIONS (Top-of-Funnel)
	// ==========================================

	// WHAT DOES YOUR COMPANY DO?
	if contains(userLower, []string{"what does your company do", "what do you do", "what do you offer", "what is marault", "marault does what", "company do", "main business", "i need service", "need service", "what service", "tell service", "describe service"}) {
		return "We help businesses and private clients make clearer, higher-quality decisions. That includes understanding revenue drivers, building forecasting systems, improving websites, and designing interfaces that are easier to trust and use. At Marault Intelligence, we combine data strategy with design excellence to solve real business problems. Whether you're struggling with unclear data, unclear positioning, or unclear systems, we help turn that into clarity."
	}

	// WHO DO YOU WORK WITH?
	if contains(userLower, []string{"who do you work with", "who is your customer", "who do you serve", "typical client", "ideal client", "target market", "who should hire you", "best fit for", "right for us", "do you work with", "companies like", "type of company", "type of client", "which industries", "what industries", "who are your clients", "who are clients", "your customer", "your clients", "is this for", "suitable for", "right fit", "match for", "good for", "companies that", "people who", "organizations that", "who benefits", "who needs this", "who would use", "would this work for", "our company", "us fit", "startup", "small business", "mid-market", "enterprise", "private equity", "real estate", "finance", "saas", "ecommerce", "retail", "non-profit"}) {
		return "We work with diverse organizations: growing companies that need better visibility into performance, leadership teams making strategic decisions, private clients managing complex financial situations, firms needing stronger positioning and digital presence, and organizations across Finance, SaaS, E-Commerce, Real Estate, Professional Services, and Healthcare. Our ideal clients value clarity, make decisions based on data, and are willing to act on insights. We typically serve companies that have reached scale where data quality impacts decisions, or individuals with complexity requiring structure. We're equally comfortable working with startups scaling fast or established enterprises optimizing performance. The common thread: you care about being clear and making confident decisions."
	}

	// COMPANY SIZE & STAGE
	if contains(userLower, []string{"do you work with startups", "startups", "early stage", "small companies", "small businesses", "enterprise", "large companies", "big companies", "mid-size", "mid-market"}) {
		return "We work across all company stages. Startups often need better decision infrastructure as they scale. Mid-market companies usually have data challenges from growth. Established enterprises often need modernization. What matters isn't size—it's whether you care about clarity and are willing to invest in systems that work. Some startups are perfect fits. Some enterprises are perfect fits. Size matters less than mindset."
	}

	// SPECIFIC INDUSTRIES
	if contains(userLower, []string{"finance", "financial services", "banking", "healthcare", "retail", "ecommerce", "saas", "technology", "real estate", "professional services", "nonprofits", "nonprofits", "consulting", "agency", "accounting", "legal", "insurance"}) {
		return "We've worked across Finance, Banking, Healthcare, E-Commerce, Retail, SaaS, Real Estate, Professional Services, Consulting, Technology, Nonprofits, Accounting, and more. Each industry has unique data challenges. What they all share: leadership teams that want to make decisions faster with more confidence. Visit our services page to see how your specific challenge maps to our solutions, or reach out to discuss your industry."
	}

	// SPECIFIC ROLES/LEADERSHIP
	if contains(userLower, []string{"ceo", "cfo", "coo", "cto", "founder", "executive", "founder", "entrepreneur", "leadership", "c-level", "c level", "leadership team", "management"}) {
		return "We primarily work with leadership teams and decision-makers: founders, executives (CEO, CFO, COO), and management teams that own strategic decisions. We also work extensively with data leaders who need analytical infrastructure support, and with private clients (high-net-worth individuals) managing complex finances. If you make decisions or influence how decisions get made, we can help."
	}

	// WHAT MAKES YOU DIFFERENT?
	if contains(userLower, []string{"what makes you different", "how are you different", "why choose you", "differentiate", "unique", "stand out", "compared to other"}) {
		return "Most firms either focus on visuals (design agencies) or focus on data (analytics firms). We do both—but more importantly, we focus on decision clarity: what's happening, why it's happening, what to do next. At Marault Intelligence, we don't believe in separation between data and design. Both should serve the same goal: reduce uncertainty and increase confidence in decisions. That's what makes us different."
	}

	// WHAT PROBLEMS DO YOU SOLVE?
	if contains(userLower, []string{"what problems do you solve", "problems solve", "what problems", "what challenges", "solve what"}) {
		return "We solve the most common challenges we see: 'We have data but no clear insight,' 'We don't know what's driving growth,' 'Our website doesn't convert or explain us well,' 'Our product feels harder to use than it should,' and 'We're making decisions without confidence.' At Marault Intelligence, we specialize in structural clarity—not cosmetic fixes. If your organization faces any of these, you're in good shape for what we do."
	}

	// HOW DO I KNOW IF I NEED THIS?
	if contains(userLower, []string{"how do i know if i need", "do i need this", "is this for us", "right for us", "fit", "need your help", "apply to us"}) {
		return "If you've ever thought 'this should be clearer,' 'we're guessing too much,' or 'our site/product isn't reflecting our quality,' then you're likely a fit. At Marault Intelligence, we work with companies and individuals that recognize clarity is competitive advantage. The best indicator is: do you have decisions that feel uncertain, or systems that feel confusing? If yes, let's talk."
	}

	// ==========================================
	// MID-LEVEL BUYER QUESTIONS
	// ==========================================

	// WHAT'S YOUR PROCESS LIKE?
	if contains(userLower, []string{"what's your process", "process", "how do you work", "approach", "methodology", "work process"}) {
		return "We don't jump straight into execution. We typically: clarify the problem, identify what's already working, structure the system (data, messaging, UX, etc.), refine and implement, and create something that stays usable over time. At Marault Intelligence, we believe the best solutions come from understanding the problem deeply first. Most engagements start with a discovery phase where we listen more than we talk."
	}

	// DO YOU CUSTOMIZE EVERYTHING?
	if contains(userLower, []string{"do you customize", "customize", "tailored", "custom work", "one size fits all"}) {
		return "Yes—but intentionally. Some engagements are fully custom (like websites or analytics systems). Others use structured foundations (like template builds) but are tailored where it matters most: messaging, hierarchy, and decision flow. At Marault Intelligence, we don't believe in unnecessary customization, but we customize everything that impacts your business outcomes. This balance keeps projects efficient while staying effective."
	}

	// HOW LONG DO PROJECTS TAKE?
	if contains(userLower, []string{"how long do projects take", "project duration", "project timeline", "how long", "weeks months"}) {
		return "Every project is unique based on scope and complexity. We prioritize quality and clarity over rushing. We'll provide exact timelines after understanding your specific needs. Visit the inquire page to discuss your project scope and get an accurate estimate."
	}

	// WHAT DO CLIENTS GET?
	if contains(userLower, []string{"what do clients get", "deliverables", "what will we get", "what do we receive", "output"}) {
		return "Not just deliverables—they get clarity on what matters, systems they can reuse, structured decision-making, and something that continues to work after we're done. At Marault Intelligence, we measure success by long-term usability, not project completion. A dashboard that sits unused is a failure. A website that doesn't convert is a failure. We deliver outcomes, not just outputs."
	}

	// DO YOU WORK WITH STARTUPS?
	if contains(userLower, []string{"startup", "early stage", "new company", "small company", "just starting"}) {
		return "Both startups and established companies—the key factor isn't size, it's complexity and seriousness. At Marault Intelligence, we work best with organizations that recognize clarity is important and are willing to invest in it. Some startups have highly complex data situations that need our help. Some large companies have simpler needs. We adapt to the situation."
	}

	// ==========================================
	// HIGH-LEVEL / EXECUTIVE QUESTIONS
	// ==========================================

	// HOW DO YOU MEASURE SUCCESS?
	if contains(userLower, []string{"measure success", "success metrics", "how do you know", "success definition"}) {
		return "Depends on the engagement, but generally: improved decision speed and confidence, clearer understanding of drivers and trade-offs, stronger conversion or user behavior, and reduced confusion internally or externally. At Marault Intelligence, we measure success by whether you're making better decisions faster. If a year after we finish working together, you're still using what we built and it's driving real value, that's success."
	}

	// WHY NOT HIRE INTERNALLY?
	if contains(userLower, []string{"why not hire internally", "hire internal team", "internal hire", "in-house solution", "do it ourselves"}) {
		return "You can—but internal teams often: are too close to the problem, lack structured frameworks, are focused on execution instead of clarity. At Marault Intelligence, we come in with: outside perspective, structured thinking, and focus on decision quality. We're also faster because this is all we do. Plus, an external perspective carries credibility with leadership and removes politics from decisions."
	}

	// HOW DO YOU PREVENT OVERENGINEERING?
	if contains(userLower, []string{"overengineered", "complex", "overly complicated", "too much", "bloat"}) {
		return "We actively avoid complexity. Our principle: if it doesn't improve decision-making, it doesn't belong. At Marault Intelligence, we've seen countless dashboards, websites, and systems that nobody uses because they're over-engineered. We ruthlessly prioritize simplicity and focus. Every feature, metric, and design element has to earn its place."
	}

	// PHILOSOPHY ON DESIGN AND ANALYTICS
	if contains(userLower, []string{"philosophy", "design philosophy", "analytics philosophy", "believe", "values"}) {
		return "Both design and analytics should serve the same goal: reduce uncertainty and increase confidence in decisions. Not decoration, not dashboards for the sake of dashboards, not unnecessary features. At Marault Intelligence, we believe clarity is competitive advantage. Every design choice, every metric, every system should make decisions easier and faster."
	}

	// DO YOU INTEGRATE WITH EXISTING TEAMS?
	if contains(userLower, []string{"integrate", "work with existing", "our team", "collaborate", "existing vendors"}) {
		return "Yes. We often: support leadership directly, collaborate with internal teams, and complement existing vendors. At Marault Intelligence, we don't believe we need to control everything—we work best as partners. We integrate with your systems, train your team, and hand off with documentation so you can maintain what we build."
	}

	// ==========================================
	// SKEPTICAL / DEFENSIVE QUESTIONS
	// ==========================================

	// WHY SHOULD WE TRUST YOU?
	if contains(userLower, []string{"why should we trust", "trust you", "credible", "proven", "track record", "how do we know"}) {
		return "You shouldn't blindly. What you should evaluate: does the thinking feel structured and clear? Do we simplify complexity instead of adding to it? Do we focus on outcomes or just deliverables? At Marault Intelligence, we aim to prove value through clarity—not claims. We're comfortable with skepticism because our work speaks for itself. Start with a conversation: visit inquire page and see if our thinking resonates."
	}

	// THIS SOUNDS VAGUE - WHAT DO YOU DELIVER?
	if contains(userLower, []string{"sounds vague", "vague", "specific deliverables", "tangible", "concrete output"}) {
		return "Fair question. We deliver tangible outputs: dashboards and analytics systems, forecasting and decision models, website design and builds, UX/UI systems and prototypes. But the real value is how those outputs improve decisions. At Marault Intelligence, we don't sell concepts—we sell working systems. You'll have something usable, maintained, and valuable long after we're done."
	}

	// WHY NOT USE A CHEAPER AGENCY?
	if contains(userLower, []string{"cheaper agency", "less expensive", "too expensive", "cost", "cheaper option", "price comparison"}) {
		return "You can. The trade-off is usually: surface-level execution vs structured thinking, output vs outcome, speed vs quality of decision-making. At Marault Intelligence, we focus on: clarity, structure, and long-term usability. Cheaper often means less strategic, less durable, and ultimately more costly because you'll rebuild it. We charge fairly for work that lasts."
	}

	// IS THIS JUST CONSULTING FLUFF?
	if contains(userLower, []string{"just consulting", "consulting fluff", "just advice", "all talk", "no execution"}) {
		return "No. We are intentionally: practical, structured, and output-driven. Every engagement results in something usable: a system, a model, a product, a site. At Marault Intelligence, we don't do 100-page reports that sit on a shelf. We deliver working infrastructure that your team uses every day."
	}

	// YOU'RE NEWER - WHY TAKE YOU SERIOUSLY?
	if contains(userLower, []string{"newer company", "new firm", "how long", "established", "proven"}) {
		return "Because seriousness comes from: quality of thinking, clarity of execution, and discipline in approach—not just age. At Marault Intelligence, we focus on delivering work that holds up under scrutiny. Our team brings deep experience in analytics, strategy, and design. If you're looking for a newer, nimbler firm with expert thinking, that's us. Let's talk."
	}

	// WE COULD DO THIS OURSELVES
	if contains(userLower, []string{"we could do this ourselves", "do it ourselves", "figure it out", "build it internally"}) {
		return "You probably could—given enough time, resources, and expertise. At Marault Intelligence, what we provide is: speed, structure, clarity, and fewer mistakes. We've done this hundreds of times. We know the patterns, the pitfalls, and the solutions. Your team can focus on running your business while we build the infrastructure. Time is your biggest cost."
	}

	// ==========================================
	// TECHNICAL / ANALYTICAL QUESTIONS
	// ==========================================

	// HOW DO YOU APPROACH MODELING?
	if contains(userLower, []string{"modeling approach", "analytics rigor", "technical approach", "model approach"}) {
		return "We prioritize: clear assumptions, driver-based frameworks, auditability (what changed and why), and scenario-based thinking over single outputs. At Marault Intelligence, we believe models should be understandable—not black boxes. Your team should be able to use them confidently and adjust them as business conditions change."
	}

	// HOW DO YOU HANDLE UNCERTAINTY?
	if contains(userLower, []string{"handle uncertainty", "uncertainty", "uncertainty forecasting", "precision"}) {
		return "We avoid false precision. Instead: scenario ranges, sensitivity analysis, decision thresholds. At Marault Intelligence, we focus on decision-readiness, not point accuracy. Forecasts that show ranges and highlight assumptions are more honest and more useful than precise predictions that ignore uncertainty."
	}

	// HOW DO YOU PREVENT OVERFITTING?
	if contains(userLower, []string{"overfitting", "fragile systems", "brittle", "robust"}) {
		return "By: keeping models interpretable, avoiding unnecessary complexity, aligning models to real business drivers. At Marault Intelligence, we build systems that survive changing business conditions. Overfitted models break when the business changes. Robust models adapt."
	}

	// UX/UI AT SYSTEMS LEVEL
	if contains(userLower, []string{"ux at system level", "ux systems", "design systems", "ui consistency"}) {
		return "We focus on: hierarchy, flow efficiency, consistency, and cognitive load reduction. At Marault Intelligence, we believe UX design should make people faster and more confident, not slower or confused. Every click, every color, every layout choice serves that goal."
	}

	// AESTHETICS VS USABILITY
	if contains(userLower, []string{"aesthetics vs usability", "beautiful vs functional", "beauty vs use"}) {
		return "Usability comes first. But: strong aesthetics reinforce trust, discipline beats decoration. At Marault Intelligence, we believe beautiful design that's hard to use is actually just bad design. And boring design that works is just incomplete. We aim for both: beautiful AND usable."
	}

	// ==========================================
	// SALES-DRIVEN / CONVERSION QUESTIONS
	// ==========================================

	// WHAT HAPPENS AFTER INQUIRY?
	if contains(userLower, []string{"what happens after inquiry", "after i inquire", "next after form", "inquiry process"}) {
		return "We: understand your situation, clarify the real problem, determine if we're a fit, and outline a clear next step. At Marault Intelligence, we don't assume every company is a fit—we'd rather have honest conversations. If it makes sense to work together, we'll propose a clear engagement path and timeline. You'll hear from us within 24 hours of inquiry."
	}

	// HOW DO YOU PRICE?
	if contains(userLower, []string{"how do you price", "pricing model", "pricing structure", "price based on"}) {
		return "Depends on: scope, complexity, and impact. We don't use one-size pricing. At Marault Intelligence, we customize pricing for each situation. Some engagements are smaller and faster, others are larger. We'll discuss options during discovery. Visit inquire page to start the conversation."
	}

	// WHAT'S THE RISK OF DOING NOTHING?
	if contains(userLower, []string{"risk of doing nothing", "wait", "do nothing", "delay", "not act"}) {
		return "Usually: continued confusion, slower decisions, missed opportunities, weaker positioning. At Marault Intelligence, we often see companies that waited too long and lost competitive advantage to better-positioned competitors or made decisions from unreliable data. The cost of inaction is often higher than the cost of action. Ready to talk? Visit inquire page."
	}

	// WHICH CLIENTS GET MOST VALUE?
	if contains(userLower, []string{"which clients get", "best fit", "ideal fit", "right fit"}) {
		return "Clients who: value clarity, are making real decisions, are willing to act on insights. At Marault Intelligence, we work best with organizations that see clarity as competitive advantage and are serious about investing in it. If you recognize a problem and are ready to solve it, you're the right fit."
	}

	// ==========================================
	// CHATBOT DEFENSE RESPONSES (Short/Powerful)
	// ==========================================

	// THIS SEEMS EXPENSIVE
	if contains(userLower, []string{"seems expensive", "too expensive", "costly", "high price"}) {
		return "The cost usually reflects the impact—better decisions tend to compound. At Marault Intelligence, we've seen companies save hundreds of thousands in avoided mistakes or gained revenue through better decisions. The question isn't 'is this expensive?'—it's 'what's the cost of unclear decisions?' Usually that's higher."
	}

	// NOT SURE IF NECESSARY
	if contains(userLower, []string{"not sure if necessary", "not sure needed", "uncertain if need"}) {
		return "That's fair. The best use case is when something already feels unclear or harder than it should be. At Marault Intelligence, we work with companies that recognize the problem. If your data feels messy, your website doesn't convert, or your decisions feel uncertain, then it's necessary. Otherwise, maybe wait until it does."
	}

	// TRIED SOMETHING LIKE THIS BEFORE
	if contains(userLower, []string{"tried before", "done something like this", "previous experience", "didn't work"}) {
		return "That happens. The difference is usually whether the output was actually usable and decision-focused. At Marault Intelligence, we focus on long-term utility, not project completion. If your previous project sits unused, that was the problem. We build systems you'll actually use."
	}

	// WE JUST NEED SOMETHING SIMPLE
	if contains(userLower, []string{"need something simple", "keep it simple", "simple solution", "just simple"}) {
		return "That's exactly the goal—but simple done well requires structure. At Marault Intelligence, we believe the best solutions look simple from the outside but have rigorous thinking underneath. 'Simple' without structure is fragile. We make it simple AND robust."
	}

	// ==========================================
	// GRAMMAR VARIATIONS & TYPOS
	// ==========================================

	// Common typos and variations
	if contains(userLower, []string{"waht", "teh", "wat", "ur", "ur services", "wut", "ur company", "ur mission"}) {
		return "I think I caught a typo there! Could you rephrase that? At Marault Intelligence, we're here to help. Feel free to ask about our services, team, pricing, timeline, or anything else. What would you like to know?"
	}

	// Default response - when no handler matches
	return "That's a great question! I may not have the specific answer in my database, but our team definitely does. The best way to get a complete answer is to visit our inquire page and share your question directly with our team—we'll respond within 24 hours with exactly what you need. In the meantime, feel free to ask about our services, qualifications, or how we work!"
}

func contains(text string, keywords []string) bool {
	for _, keyword := range keywords {
		// Exact match
		if strings.Contains(text, keyword) {
			return true
		}
		// Fuzzy match for typos (handles common misspellings)
		if fuzzyMatch(text, keyword) {
			return true
		}
	}
	return false
}

// levenshteinDistance calculates the minimum number of edits needed to transform one string into another
func levenshteinDistance(s1, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)
	len1 := len(r1)
	len2 := len(r2)

	if len1 == 0 {
		return len2
	}
	if len2 == 0 {
		return len1
	}

	// Create a matrix to store distances
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1, // deletion
				min(matrix[i][j-1]+1, // insertion
					matrix[i-1][j-1]+cost), // substitution
			)
		}
	}

	return matrix[len1][len2]
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// fuzzyMatch checks if text contains a keyword with tolerance for typos
// Returns true if the keyword is found with a similarity threshold
func fuzzyMatch(text string, keyword string) bool {
	// Minimum word length to fuzzy match
	if len(keyword) < 4 {
		return false // Don't fuzzy match very short words
	}

	// Split text into words and check each word
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'))
	})

	for _, word := range words {
		// Calculate similarity
		distance := levenshteinDistance(word, keyword)
		maxLen := len(word)
		if len(keyword) > maxLen {
			maxLen = len(keyword)
		}

		// Allow up to 1/3 of the word length as edits for longer words
		// For example: "dashbord" vs "dashboard" (1 edit, 9 chars = 11% difference)
		threshold := (maxLen + 2) / 3 // Allows ~33% character difference

		if distance <= threshold {
			return true
		}
	}

	return false
}



func main() {
    mux := http.NewServeMux()
    /* =========================
       Static files
    ========================= */
    fileServer := http.FileServer(http.Dir("./static"))
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))
    /* =========================
       Debug
    ========================= */
    mux.HandleFunc("/debug-ua", func(w http.ResponseWriter, r *http.Request) {
        ua := r.UserAgent()
        tmpl := getBaseTemplate(r)
        fmt.Fprintf(w, "User-Agent: %s\n\nTemplate: %s", ua, tmpl)
    })
    /* =========================
       Core pages
    ========================= */
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/approach", approachHandler)
    mux.HandleFunc("/executive-team", executiveTeamHandler)
    mux.HandleFunc("/contact", contactHandler)
    mux.HandleFunc("/inquire", inquireHandler)
    mux.HandleFunc("/philosophy", philosophyHandler)

	/* =========================
	   API endpoints
	========================= */
	mux.HandleFunc("/api/chat", chatHandler)

	/* =========================
	   Services
	========================= */
	mux.HandleFunc("/services", servicesHandler)

	mux.HandleFunc(
		"/services/data-visibility-audit",
		servicePageHandler(
			"data-visibility-audit.html",
			"Data Visibility Audit | Marault Intelligence",
		),
	)

	mux.HandleFunc(
		"/services/revenue-customer-analytics",
		servicePageHandler(
			"revenue.html",
			"Revenue & Customer Analytics | Marault Intelligence",
		),
	)

	mux.HandleFunc(
		"/services/custom-website-build",
		servicePageHandler(
			"custom-website-build.html",
			"Custom Website Build | Marault Intelligence",
		),
	)

	mux.HandleFunc(
		"/services/executive-dashboards-reporting",
		servicePageHandler(
			"executive-dashboards-reporting.html",
			"Executive Dashboards & Reporting | Marault Intelligence",
		),
	)

	mux.HandleFunc(
		"/services/forecasting-decision-modeling",
		servicePageHandler(
			"forecasting-decision-modeling.html",
			"Forecasting & Decision Modeling | Marault Intelligence",
		),
	)

	mux.HandleFunc(
		"/services/private-client-analytics",
		servicePageHandler(
			"private-client-analytics.html",
			"Private Client Analytics | Marault Intelligence",
		),
	)

	mux.HandleFunc(
		"/services/template-based-build",
		servicePageHandler(
			"template-based-build.html",
			"Template-Based Build | Marault Intelligence",
		),
	)

	mux.HandleFunc(
		"/services/website-redesign",
		servicePageHandler(
			"website-redesign.html",
			"Website Redesign | Marault Intelligence",
		),
	)

	mux.HandleFunc(
		"/services/ux-ui-design",
		servicePageHandler(
			"ux-ui-design.html",
			"UX/UI Design | Marault Intelligence",
		),
	)

	port := os.Getenv("PORT")
	if port == "" {
    port = "8080"
	}
	log.Println("Starting server on 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, mux))
}

/*
	=========================
	  HOME PAGE

=========================
*/
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(
		getBaseTemplate(r),
		"./internal/templates/home.html",
		"./internal/templates/chatbot.html",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Page  string
	}{
		Title: "Marault Intelligence",
		Page:  "home",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

/*
	=========================
	  APPROACH PAGE

=========================
*/
func approachHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		getBaseTemplate(r),
		"./internal/templates/approach.html",
		"./internal/templates/chatbot.html",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Page  string
	}{
		Title: "The Marault Approach | Marault Intelligence",
		Page:  "approach",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

/*
	=========================
	  EXECUTIVE TEAM

=========================
*/
func executiveTeamHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		getBaseTemplate(r),
		"./internal/templates/executive.html",
		"./internal/templates/chatbot.html",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Page  string
	}{
		Title: "Executive Team | Marault Intelligence",
		Page:  "executive",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

/*
	=========================
	  SERVICES OVERVIEW

=========================
*/
func servicesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		getBaseTemplate(r),
		"./internal/templates/services.html",
		"./internal/templates/chatbot.html",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Page  string
	}{
		Title: "Services | Marault Intelligence",
		Page:  "services",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

/*
	=========================
	  CONTACT

=========================
*/
func contactHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		getBaseTemplate(r),
		"./internal/templates/contact.html",
		"./internal/templates/chatbot.html",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Page  string
	}{
		Title: "Contact | Marault Intelligence",
		Page:  "contact",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

/*
	=========================
	  INQUIRE

=========================
*/
func inquireHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(
			getBaseTemplate(r),
			"./internal/templates/inquire.html",
			"./internal/templates/chatbot.html",
		)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
			Page  string
		}{
			Title: "Inquire | Marault Intelligence",
			Page:  "inquire",
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println(err)
		}
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		company := r.FormValue("company")
		message := r.FormValue("message")
		services := r.Form["services"]

		selectedServices := strings.Join(services, ", ")

		err := sendInquiryEmail(name, email, company, selectedServices, message)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to send message", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles(
			getBaseTemplate(r),
			"./internal/templates/thankyou.html",
			"./internal/templates/chatbot.html",
		)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
			Page  string
		}{
			Title: "Thank You | Marault Intelligence",
			Page:  "thankyou",
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println(err)
		}
	}
}

/*
	=========================
	  EMAIL SENDER

=========================
*/
func sendInquiryEmail(name, email, company, services, message string) error {
	from := "caroline@maraultintelligence.com"
	password := "ljblmutreemgdd"

	to := []string{
		"caroline@maraultintelligence.com",
		"lindsey@maraultintelligence.com",
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	body := fmt.Sprintf(
		"New Inquiry\n\nName: %s\nEmail: %s\nCompany: %s\nServices: %s\n\nMessage:\n%s",
		name, email, company, services, message,
	)

	msg := "From: " + from + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: New Website Inquiry\n\n" + body

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
}


func philosophyHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		getBaseTemplate(r),
		"./internal/templates/philosophy.html",
		"./internal/templates/chatbot.html",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Page  string
	}{
		Title: "Data Philosophy | Marault Intelligence",
		Page:  "philosophy",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}
