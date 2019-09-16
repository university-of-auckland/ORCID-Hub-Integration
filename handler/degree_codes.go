package main

var degreeCodes = map[string]string{
	"BA":                   "Bachelor of Arts",
	"BA(HONS)":             "Bachelor of Arts with Honours",
	"BA(JAP)":              "Bachelor of Arts (Japanese)",
	"BA(SOCSC)":            "Bachelor of Arts (Social Sciences)",
	"BACC":                 "Bachelor of Accountancy",
	"BAGR":                 "Bachelor of Agriculture",
	"BAGRECON":             "Bachelor of Agricultural Economics",
	"BAGRECON(HONS)":       "Bachelor of Agricultural Economics with Honours",
	"BAGRSC":               "Bachelor of Agricultural Science",
	"BAGRSC(HONS)":         "Bachelor of Agricultural Science (Honours)",
	"BAHED":                "Bachelor of Adult and Higher Education",
	"BAPPLCOMP":            "Bachelor of Applied Computing",
	"BAPPLECON":            "Bachelor of Applied Economics",
	"BAPPLECON(HONS)":      "Bachelor of Applied Economics with Honours",
	"BAPPLSC":              "Bachelor of Applied Science",
	"BARCH":                "Bachelor of Architecture",
	"BARCH(HONS)":          "Bachelor of Architecture (Honours)",
	"BARTDES":              "Bachelor of Art and Design",
	"BARTDES(HONS)":        "Bachelor of Art and Design with Honours",
	"BAS":                  "Bachelor of Architectural Science",
	"BASC":                 "Bachelor of Agricultural Science",
	"BAV":                  "Bachelor of Aviation",
	"BAV(HONS)":            "Bachelor of Aviation with Honours",
	"BAVMAN":               "Bachelor of Aviation Management",
	"BAVMAN(HONS)":         "Bachelor of Aviation Management with Honours",
	"BBIM":                 "Bachelor of Business and Information Management",
	"BBIOMEDSC":            "Bachelor of Biomedical Science",
	"BBIS":                 "Bachelor of Business Information Systems",
	"BBMEDSC(HONS)":        "Bachelor of Biomedical Science with Honours",
	"BBS":                  "Bachelor of Business Studies",
	"BBS(HONS)":            "Bachelor of Business Studies with Honours",
	"BBSC":                 "Bachelor of Building Science",
	"BBSC(HONS)":           "Bachelor of Building Science with Honours",
	"BBUS":                 "Bachelor of Business",
	"BBUS(HONS)":           "Bachelor of Business (Honours)",
	"BBUSINF":              "Bachelor of Business Information",
	"BC":                   "Bachelor of Communications",
	"BC(QS)":               "Bachelor of Construction (Quantity Surveying)",
	"BCA":                  "Bachelor of Commerce and Administration",
	"BCA(HONS)":            "Bachelor of Commerce and Administration with Honours",
	"BCAPSC":               "Bachelor of Consumer and Applied Sciences",
	"BCAPSC(HONS)":         "Bachelor of Consumer and Applied Sciences with Honours",
	"BCM":                  "Bachelor of Commerce and Management",
	"BCMS":                 "Bachelor of Computing and Mathematical Sciences",
	"BCOM":                 "Bachelor of Commerce",
	"BCOM(AG)":             "Bachelor of Commerce (Agricultural) L [now Agriculture]",
	"BCOM(AG)(HONS)":       "Bachelor of Commerce (Agricultural) with Honours",
	"BCOM(FORESTRY)":       "Bachelor of Commerce- Forestry",
	"BCOM(H&IM)":           "Bachelor of Commerce (Hotel and Institutional Management)",
	"BCOM(HONS)":           "Bachelor of Commerce (Honours)",
	"BCOM(HORT)":           "Bachelor of Commerce (Horticultural)",
	"BCOM(M&TM)":           "Bachelor of Commerce (Manufacturing and Technology Management)",
	"BCOM(T&L)":            "Bachelor of Commerce (Transport and Logistics)",
	"BCOM(TOURISM)":        "Bachelor of Commerce (Tourism)",
	"BCOM(TRANSPORT)":      "Bachelor of Commerce (Transport)",
	"BCOM(VPM)":            "Bachelor of Commerce (Valuation and Property Management)",
	"BCS":                  "Bachelor of Communications Studies",
	"BD":                   "Bachelor of Divinity",
	"BDANCE":               "Bachelor of Dance",
	"BDEFSTUDS":            "Bachelor of Defence Studies",
	"BDENTTECH":            "Bachelor of Dental Technology",
	"BDES":                 "Bachelor of Design",
	"BDS":                  "Bachelor of Dental Surgery",
	"BE":                   "Bachelor of Engineering",
	"BE(HONS)":             "Bachelor of Engineering with Honours",
	"BED":                  "Bachelor of Education",
	"BED(ADULTED)":         "Bachelor of Education (Adult Education)",
	"BED(EARLY":            "Childhood Teaching) Bachelor of Education (Early Childhood Teaching)",
	"BED(HONS)":            "Bachelor of Education with Honours",
	"BED(SECTCHG)":         "Bachelor of Education (Secondary Teaching)",
	"BED(TCHG)":            "Bachelor of Education (Teaching)",
	"BED(TCHG)HONS":        "Bachelor of Education (Teaching) with Honours",
	"BED(TESOL)":           "Bachelor of Education (Teaching English to Speakers of Other Languages)",
	"BEDSC":                "Bachelor of Education in Science",
	"BEM":                  "Bachelor of Environmental Management",
	"BEM(HONS)":            "Bachelor of Environmental Management with Honours",
	"BENG(HONS)":           "Bachelor of Engineering (Honours)",
	"BENGTECH":             "Bachelor of Engineering Technology",
	"BFA":                  "Bachelor of Fine Arts",
	"BFA(HONS)":            "Bachelor of Fine Arts with Honours",
	"BFORSC":               "Bachelor of Forestry Science",
	"BFSC":                 "Bachelor of Forestry Science",
	"BGD":                  "Bachelor of Graphic Design",
	"BHB":                  "Bachelor of Human Biology",
	"BHEALSC":              "Bachelor of Health Science",
	"BHLTHSC":              "Bachelor of Health Sciences",
	"BHM":                  "Bachelor of Hospitality Management",
	"BHORT":                "Bachelor of Horticulture",
	"BHORTSC":              "Bachelor of Horticultural Science",
	"BHORTSC(HONS)":        "Bachelor of Horticultural Science with Honours",
	"BHSC":                 "Bachelor of Health Science",
	"BHSC(HONS)":           "Bachelor of Health Science with Honours",
	"BHSC(MIDWIFERY)":      "Bachelor of Health Science (Midwifery)",
	"BHSC(NURSING)":        "Bachelor of Health Science (Nursing)",
	"BHSC(OCCUPATIONAL":    "Therapy) Bachelor of Health Science (Occupational Therapy)",
	"BHSC(PHYSIOTHERAPY))": "Bachelor of Health Science (Physiotherapy)",
	"BIHM":                 "Bachelor of International Hospitality Management",
	"BIT":                  "Bachelor of Information Technology",
	"BIT(HONS":             "Bachelor of Information Technology with Honours",
	"BL&A":                 "Bachelor of Laws and Arts",
	"BL&CMS":               "Bachelor of Laws and Computing and Mathematical Sciences",
	"BL&MS":                "Bachelor of Laws and Management Studies",
	"BL&SC":                "Bachelor of Laws and Science",
	"BL&SC(TECH)":          "Bachelor of Laws and Science (Technology)",
	"BL&SOCSC":             "Bachelor of Laws and Social Sciences",
	"BLA":                  "Bachelor of Landscape Architecture",
	"BLA(HONS)":            "Bachelor of Landscape Architecture with Honours",
	"BLIBS":                "Bachelor of Liberal Studies",
	"BLPA":                 "Bachelor of Legal Policy and Administration",
	"BLS":                  "Bachelor of Leisure Studies W [now BSLS]",
	"BLS(HONS)":            "Bachelor of Leisure Studies W [now BSLS(Hons)]",
	"BMAST":                "Bachelor of Maori Studies",
	"BMD":                  "Bachelor of Maori Development",
	"BMEDSC":               "Bachelor of Medical Science",
	"BMEDSC(HONS)":         "Bachelor of Medical Science with Honours",
	"BMID":                 "Bachelor of Midwifery",
	"BMLS":                 "Bachelor of Medical Laboratory Science",
	"BMLSC":                "Bachelor of Medical Laboratory Science",
	"BMPA":                 "Bachelor of Maori Performing Arts",
	"BMPD":                 "Bachelor of Maori Planning and Development",
	"BMS":                  "Bachelor of Management Studies",
	//"BMS":                  "Bachelor of Maori Studies",
	"BMTRADARTS": "Bachelor of Maori Traditional Arts",
	"BMUS":       "Bachelor of Music",
	//"BMUS":                 "Executant Bachelor of Music Executant",
	"BMUS(HONS)":       "Bachelor of Music (Honours)",
	"BMUS(PERF)":       "Bachelor of Music (Performance)",
	"BMUS(PERF)(HONS)": "Bachelor of Music (Performance) (Honours)",
	"BMUSED":           "Bachelor of Music Education",
	"BMVA":             "Bachelor of Maori Visual Arts",
	"BN":               "Bachelor of Nursing",
	"BNURS":            "Bachelor of Nursing",
	"BNURS(HONS)":      "Bachelor of Nursing with Honours",
	"BOPTOM":           "Bachelor of Optometry",
	"BP&RMGT":          "Bachelor of Parks and Recreation Management",
	"BP&RMGT(HONS)":    "Bachelor of Parks and Recreation Management with Honours",
	"BPA":              "Bachelor of Property Administration",
	"BPERARTS":         "Bachelor of Performing Arts",
	"BPERFDES":         "Bachelor of Performance Design",
	"BPHARM":           "Bachelor of Pharmacy",
	"BPHARM(HONS)":     "Bachelor of Pharmacy with Honours",
	"BPHED":            "Bachelor of Physical Education",
	"BPHED(HONS)":      "Bachelor of Physical Education with Honours",
	"BPHIL":            "Bachelor of Philosophy",
	"BPHTY":            "Bachelor of Physiotherapy",
	"BPLAN":            "Bachelor of Planning",
	"BPR&TM":           "Bachelor of Park, Recreation and Tourism Management",
	"BPROP":            "Bachelor of Property",
	"BPROP(HONS)":      "Bachelor of Property (Honours)",
	"BRM":              "Bachelor of Recreation Management",
	"BRP":              "Bachelor of Resource and Environmental Planning",
	"BRP(HONS)":        "Bachelor of Resource and Environmental Planning with Honours",
	"BRS":              "Bachelor of Resource Studies L [now Environmental Management]",
	"BSC":              "Bachelor of Science",
	"BSC(HONS)":        "Bachelor of Science (Honours)",
	"BSC(TECH)":        "Bachelor of Science (Technology)",
	"BSCED":            "Bachelor of Science Education",
	"BSD":              "Bachelor of Spatial Design",
	"BSLS":             "Bachelor of Sport Leisure Studies",
	"BSLS(HONS)":       "Bachelor of Sport and Leisure Studies",
	"BSLT":             "Bachelor of Speech and Language Therapy",
	"BSOCSC":           "Bachelor of Social Science",
	// "BSOCSC":               "Bachelor of Social Sciences",
	"BSOCSC(HONS)":  "Bachelor of Social Sciences with Honours",
	"BSPCHLANGTHER": "Bachelor of Speech and Language Therapy",
	"BSPTSTUDS":     "Bachelor of Sports Studies",
	"BSR":           "Bachelor of Sport and Recreation",
	"BSURV":         "Bachelor of Surveying",
	"BSURV(HONS)":   "Bachelor of Surveying with Honours",
	"BSW":           "Bachelor of Social Work",
	"BTCHG":         "Bachelor of Teaching",
	"BTCHG(HONS)":   "Bachelor of Teaching with Honours",
	"BTCHG(PRIM)":   "Bachelor of Teaching (Primary)",
	"BTCHG(SEC)":    "Bachelor of Teaching (Secondary)",
	"BTEACH":        "Bachelor of Teaching",
	"BTECH":         "Bachelor of Technology",
	"BTHEOL":        "(Hons) Bachelor of Theology with Honours",
	// "BTHEOL":               "Bachelor of Theology",
	"BTOUR":         "Bachelor of Tourism",
	"BTOURMGT(HON)": "Bachelor of Tourism Management with Honours",
	"BTP":           "Bachelor of Town Planning",
	"BTSM":          "Bachelor of Tourism and Services Management",
	"BV&O":          "Bachelor of Viticulture and Oenology",
	"BVA":           "Bachelor of Visual Arts",
	"BVSC":          "Bachelor of Veterinary Science",
	"D(UOA)":        "Doctor of the University (of Auckland)",
	"DBA":           "Doctor of Education",
	"DCLINPSY":      "Doctor of Clinical Psychology",
	"DCOM":          "Doctor of Commerce",
	"DDS":           "Doctor of Dental Science",
	"DDSC":          "Doctor of Dental Science",
	"DENG":          "Doctor of Engineering",
	"DFA":           "Doctor of Fine Arts",
	"DHSC":          "Doctor of Health Science",
	"DJUR":          "Doctor of Jurisprudence",
	"DLIT":          "Doctor of Literature",
	"DLITT":         "Doctor of Literature",
	"DMA":           "Doctor of Musical Arts",
	"DMID":          "Doctor of Midwifery",
	"DMUS":          "Doctor of Music",
	"DNATRES":       "Doctor of Natural Resources",
	"DNSG":          "Doctor of Nursing",
	"DOCFA":         "Doctor of Fine Arts",
	"DPHARM":        "Doctor of Pharmacy",
	"DPHIL":         "Doctor of Philosophy W [now PhD]",
	"DPHYS":         "Doctor of Physiotherapy",
	"DSC":           "Doctor of Science",
	"DU":            "Doctor of the University",
	"EDD":           "Doctor of Education",
	"HOND":          "Doctor of the University",
	"IMBA":          "International Master of Business Administration",
	"LITD":          "Doctor of Literature",
	"LITTD":         "Doctor of Letters",
	// "LITTD":                "Doctor of Literature",
	"LLB":              "Bachelor of Laws",
	"LLB(HONS)":        "Bachelor of Laws (Honours)",
	"LLD":              "Doctor of Laws",
	"LLM":              "Master of Laws",
	"LLM(INTLAW&POLS)": "Master of Laws (International Law and Politics)",
	"LLMENVIR":         "Master of Laws (Environmental)",
	"MA":               "Master of Arts",
	"MA(APPLIED)":      "Master of Arts (Applied)",
	"MA(ART&DES)":      "Master of Arts (Art and Design)",
	"MA(COMMST)":       "Master of Arts (Communication Studies)",
	"MAF":              "Master of Applied Finance",
	"MAGR":             "Master of Agriculture",
	"MAGRECON":         "Master of Agricultural Economics",
	"MAGRSC":           "Master of Agricultural Science",
	"MAPA":             "Master of Asia-Pacific Affairs",
	"MAPPLECON":        "Master of Applied Economics",
	"MAPPLPSYCH":       "Master of Applied Psychology",
	"MAPPLSC":          "Master of Applied Science",
	"MAPPLSTAT":        "Master of Applied Statistics",
	"MARCH":            "Master of Architecture",
	"MAS":              "Master of Architectural Science",
	"MASC":             "Master of Agricultural Science",
	"MAUD":             "Master of Audiology",
	"MAV":              "Master of Aviation",
	"MBA":              "Master of Business Administration",
	"MBBS":             "Master of Medicine and Bachelor of Surgery",
	"MBCHB":            "Master of Medicine and Bachelor of Surgery",
	"MBHL":             "Master of Bioethics and Health Law",
	"MBLDGSC":          "Master of Building Science",
	"MBMEDSC(HONS)":    "Master of Biomedical Science",
	"MBS":              "Master of Business Studies",
	"MBSC":             "Master of Building Science",
	"MBUS":             "Master of Business",
	"MBUSINF":          "Master of Business Information",
	"MCA":              "Master of Commerce and Administration",
	"MCAPSC":           "Master of Consumer and Applied Sciences",
	"MCLINPHARM":       "Master of Clinical Pharmacy",
	"MCM":              "Master of Commerce and Management",
	"MCMS":             "Master of Computing and Mathematical Sciences",
	"MCOM":             "Master of Commerce",
	"MCOM(AG)":         "Master of Commerce (Agricultural)",
	"MCOMDENT":         "Master of Community Dentistry",
	"MCOMLAW":          "Master of Commercial Law",
	"MCOMMS":           "Master of Communications",
	"MCOMPSC":          "Master of Computer Science",
	"MCONBIO":          "Master of Conservation Biology",
	"MCONSC":           "Master of Conservation Science",
	"MCOUNS":           "Master of Counselling",
	"MCPA":             "Master of Creative and Performing Arts",
	"MD":               "Doctor of Medicine",
	"MDA":              "Master of Development Administration",
	"MDAIRYSCTECH":     "Master of Dairy Science and Technology",
	"MDANCE":           "Master of Dance",
	"MDANCEST":         "Master of Dance Studies",
	"MDES":             "Master of Design",
	"MDEV":             "Stud Master of Development Studies",
	"MDS":              "Master of Dental Surgery",
	"ME":               "Master of Engineering",
	"MECOM":            "Master of Electronic Commerce",
	"MED":              "Master of Education",
	"MEDADMIN":         "Master of Educational Administration",
	"MEDMGT":           "Master of Educational Management",
	"MEDPSYCH":         "Master of Educational Psychology",
	"MEDSTUDS":         "Master of Educational Studies",
	"MEFE":             "Master of Engineering in Fire Engineering",
	"MEM":              "Master of Engineering in Management",
	"MEMGT":            "Master of Engineering Management",
	"MENGST":           "Master of Engineering Studies",
	"MENGSTUDS":        "Master of Engineering Studies",
	"MENTR":            "Master of Entrepreneurship",
	"MENVLS":           "Master of Environmental Legal Studies",
	"MENVSTUD":         "Master of Environmental Studies",
	"MEP":              "Master of Environmental Planning",
	// "MEP":                  "Master of Environmental Policy",
	"MERG":       "Master of Ergonomics",
	"MET":        "Master of Engineering in Transportation",
	"MFA":        "Master of Fine Arts",
	"MFIN":       "Math Master of Financial Mathematics",
	"MGP":        "Master of General Practice",
	"MGUIDCOUNS": "Master of Guidance and Counselling",
	"MHB":        "Master of Human Biology",
	"MHEALSC":    "Master of Health Sciences",
	"MHEALTHMGT": "Master of Health Management",
	"MHORTSC":    "Master of Horticultural Science",
	"MHSC":       "Master of Health Sciences",
	// "MHSC":                 "Master of Home Science",
	"MIHM":         "Master of International Hospitality Management",
	"MIM":          "Master of Information Management",
	"MINDS":        "Master of Indigenous Studies",
	"MINFOTECH":    "Master of Information Technology",
	"MINFSC":       "Master of Information Sciences",
	"MINTBUS":      "Master of International Business",
	"MINTLAW&POLS": "Master of International Law and Politics",
	"MINTST":       "Master of International Studies",
	"MIPD":         "Master of Indigenous Planning and Development",
	"MIR":          "Master of International Relations",
	"MIS":          "Master of Information Systems",
	"MIT":          "Master of Innovation Technology",
	"MJUR":         "Master of Jurisprudence",
	"ML&A":         "Master of Laws and Arts",
	"MLA":          "Master of Landscape Architecture",
	"MLIS":         "Master of Library and Information Studies",
	"MLITT":        "Master of Literature",
	"MLS":          "Master of Leisure Studies",
	"MMEDSC":       "Master of Medical Sciences",
	"MMGT":         "Master of Management",
	"MMID":         "Master of Midwifery",
	"MMIDW":        "Master of Midwifery",
	"MMIN":         "Master of Ministry",
	"MMLSC":        "Master of Medical Laboratory Science",
	"MMPHTY":       "Master of Manipulative Physiotherapy",
	"MMS":          "Master of Management Studies",
	"MMUS":         "Master of Music",
	"MMUSTHER":     "Master of Music Therapy",
	"MMVA":         "Master of Maori Visual Arts",
	"MN":           "Master of Nursing",
	"MN(CLINICAL)": "Master of Nursing (Clinical)",
	"MNRM&EE":      "Master of Natural Resources Management and Ecological Engineering",
	"MNURS":        "Master of Nursing",
	"MNZS":         "Master of New Zealand Studies",
	"MOPHTH":       "Master of Ophthalmology",
	"MOR":          "Master of Operations Research",
	"MP&RMGT":      "Master of Parks and Recreation Management",
	"MPA(EXEC)":    "Master of Public Administration (Executive)",
	"MPH":          "Master of Public Health",
	"MPHARM":       "Master of Pharmacy",
	"MPHARMPRAC":   "Master of Pharmacy Practice",
	"MPHC":         "Master of Primary Health Care",
	"MPHED":        "Master of Physical Education",
	"MPHIL":        "Master of Philosophy",
	"MPHIST":       "Master of Public History",
	"MPHTY":        "Master of Physiotherapy",
	"MPLAN":        "Master of Planning",
	"MPLANPRAC":    "Master of Planning Practice",
	"MPM":          "Master of Public Management",
	"MPP":          "Master of Public Policy",
	"MPR&TM":       "Master of Park, Recreation and Tourism Management",
	"MPROFSTUDS":   "Master of Professional Studies",
	"MPROP":        "Master of Property",
	"MPROPSTUDS":   "Master of Property Studies",
	"MRP":          "Master of Resource and Environmental Planning",
	"MRRP":         "Master of Regional and Resource Planning",
	"MRS":          "Master of Resource Studies",
	"MS":           "Master of Surgery",
	"MSC":          "Master of Science",
	"MSC(TECH)":    "Master of Science (Technology)",
	"MSCED":        "Master of Science Education",
	"MSLS":         "Master of Sport and Leisure Studies",
	"MSLTPRAC":     "Master of Speech Language Therapy Practice",
	"MSOCSC":       "Master of Social Sciences",
	"MSPED":        "Master of Special Education",
	"MSS":          "Master of Strategic Studies",
	"MSURV":        "Master of Surveying",
	"MSW":          "Master of Social Work",
	"MSW(APP)":     "Master of Social Work(Applied)",
	"MTA":          "Master of Theatre Arts",
	"MTAXS":        "Master of Taxation Studies",
	"MTCHG":        "Master of Teaching",
	"MTEACH":       "Master of Teaching",
	"MTECH":        "Master of Technology",
	"MTESOL":       "Master of Teaching English to Speakers of Other Languages",
	"MTH":          "Master of Theology",
	"MTHEOL":       "Master of Theology",
	"MTM":          "Master of Technology Management",
	"MTOUR":        "Master of Tourism",
	"MTOURMGT":     "Master of Tourism Management",
	"MTP":          "Master of Town Planning",
	"MUSB":         "Bachelor of Music",
	"MUSB(HONS)":   "Bachelor of Music with Honours",
	"MUSD":         "Doctor of Music",
	"MVS":          "Master of Veterinary Studies",
	"MVSC":         "Master of Veterinary Science",
	"PHD":          "Doctor of Philosophy",
}
