//
//  NotificationPublishViewController.m
//  Squirrel
//
//  Created by JamesMao on 4/12/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "NotificationPublishViewController.h"
#import "AFNetworking.h"
#import "Define.h"
#import "DataModel/MeViewData.h"
#import "URLManager/URLManager.h"
#import "OneClassRegisterStudentViewController.h"
#import "DataModel/OneClassReigsterStudentData.h"
@interface NotificationPublishViewController ()
{
    IBOutlet UITextView* classNameTextView_;
    IBOutlet UITextField* classTimeTextField_;
    IBOutlet UITextField* classTeachTextField_;
    IBOutlet UITextField* classMaxStudentNumberTextField_;
    IBOutlet UITextView* classDescriptionTextView_;
    
    IBOutlet UIButton* deleteButton_;
    IBOutlet UIButton* registerButton_;
    
    IBOutlet UILabel* registerCountLabel_;
    
    NotificationDataItem* notificationDataItem_;
    UserRegisterDataItem* userRegisterDataItem_;
    
    EditType editType_;
    
    OneClassRegisterStudentViewController* oneClassRegisterStudentViewController_;
    
    BOOL bCurrentUserRegistered_;
}
- (IBAction)onDelete:(id)sender;
- (IBAction)onRegister:(id)sender;
- (IBAction)onStudentDetail:(id)sender;

@end


@implementation NotificationPublishViewController

- (id)init
{
    if (self = [super init])
    {
        notificationDataItem_ = [[NotificationDataItem alloc] init];
        
        userRegisterDataItem_ = [[UserRegisterDataItem alloc] init];
        
        self.hidesBottomBarWhenPushed = YES;
        editType_ = EditType_Add;
        
        oneClassRegisterStudentViewController_ = [[OneClassRegisterStudentViewController alloc] init];
    }
    return self;
}

- (void)viewDidLoad
{
    [super viewDidLoad];
    
    self.navigationItem.rightBarButtonItem = [[UIBarButtonItem alloc] initWithTitle:@"Save"
                                                                              style:UIBarButtonItemStyleDone
                                                                             target:self
                                                                             action:@selector(onSave)];
    
    self.navigationController.navigationBar.titleTextAttributes = @{NSForegroundColorAttributeName: [UIColor whiteColor]};
    
    self.navigationController.navigationBar.tintColor = [UIColor whiteColor];
    self.title = @"Class Description";
    
}

- (void)onSave
{
    [notificationDataItem_ setClassName: classNameTextView_.text];
    [notificationDataItem_ setClassTime: classTimeTextField_.text];
    [notificationDataItem_ setClassTeacher: classTeachTextField_.text];
    [notificationDataItem_ setClassDescription: classDescriptionTextView_.text];
    [notificationDataItem_ setClassMaxStudent: classMaxStudentNumberTextField_.text];
    
    [[AdminPublishViewData singleton] setNotificationDataItem:notificationDataItem_];
    
    [self postClass];
}

- (void)postClass
{
    NSString* urlString = [URLManager urlString:URLTypeClassPost];
    NSDictionary *parameters = [notificationDataItem_ getDataItemEelements];
    [self request:RequestTypePost urlString:urlString parameters:parameters];
}

- (void)getOneClass
{
    NSString* key = [notificationDataItem_ getUniqueKey];
    NSString* urlString = [URLManager urlString:URLTypeOneClass variableKey:key];
    [self request:RequestTypeGet urlString:urlString parameters:nil];
}

- (void)getOneClassUserCount
{
    registerCountLabel_.text = @"";
    
    NSString* key = [notificationDataItem_ getUniqueKey];
    NSString* urlString = [URLManager urlString:URLTypeOneClassUserCount variableKey:key];
    [self request:RequestTypeGet urlString:urlString parameters:nil];
}

- (void)deleteOneClass
{
    NSString* key = [notificationDataItem_ getUniqueKey];
    NSString* urlString = [URLManager urlString:URLTypeOneClass variableKey:key];
    [self request:RequestTypeDelete urlString:urlString parameters:nil];
    [self.navigationController popViewControllerAnimated:YES];
}

- (void)registerOneUser
{
    NSString* uniqueNameKey = [MeViewData singleton].userUniqueName;
    [userRegisterDataItem_ setUniqueKey:[notificationDataItem_ getUniqueKey]];
    [userRegisterDataItem_ setClassName:[notificationDataItem_ getClassName]];
    [userRegisterDataItem_ setUserName:uniqueNameKey];
    [userRegisterDataItem_ setUserUniqueKey:uniqueNameKey];
    NSDictionary *parameters = [userRegisterDataItem_ getDataItemEelements];
    
    NSString* urlString = [URLManager urlString:URLTypeUserRegister];
    [self request:RequestTypePost urlString:urlString parameters:parameters];
}

- (void)unregisterOneUser
{
    NSString* uniqueNameKey = [MeViewData singleton].userUniqueName;
    [userRegisterDataItem_ setUniqueKey:[notificationDataItem_ getUniqueKey]];
    [userRegisterDataItem_ setClassName:[notificationDataItem_ getClassName]];
    [userRegisterDataItem_ setUserName:uniqueNameKey];
    [userRegisterDataItem_ setUserUniqueKey:uniqueNameKey];
    NSDictionary *parameters = [userRegisterDataItem_ getDataItemEelements];
    
    NSString* urlString = [URLManager urlString:URLTypeUserUnregister];
    [self request:RequestTypePost urlString:urlString parameters:parameters];

}

- (void)queryRegisterStatus
{
    NSString* uniqueNameKey = [MeViewData singleton].userUniqueName;
    [userRegisterDataItem_ setUniqueKey:[notificationDataItem_ getUniqueKey]];
    [userRegisterDataItem_ setClassName:[notificationDataItem_ getClassName]];
    [userRegisterDataItem_ setUserName:uniqueNameKey];
    [userRegisterDataItem_ setUserUniqueKey:uniqueNameKey];
    NSDictionary *parameters = [userRegisterDataItem_ getDataItemEelements];
    
    NSString* urlString = [URLManager urlString:URLTypeQueryRegisterStatus];
    [self request:RequestTypeGet urlString:urlString parameters:parameters];
}


- (void)onSuccess:(AFHTTPRequestOperation *)operation responseObject:(id)responseObject
{
    [super onSuccess:operation responseObject:responseObject];
    NSString* resultString = [responseObject objectForKey:@"result"];
    BOOL bDelete = NO;
    if ([resultString rangeOfString:@"delete"].length > 0)
    {
        bDelete = YES;
    }
    NSString* key = [notificationDataItem_ getUniqueKey];
    NSString* oneClassURLString = [URLManager urlString:URLTypeOneClass variableKey:key];
    NSString* classPostURLString = [URLManager urlString:URLTypeClassPost];
    NSString* candidate = [[[operation response] URL] absoluteString];
    if (bDelete)
    {
        NSLog(@"delete");
    }
    else if ([oneClassURLString isEqualToString:candidate])
    {
        [notificationDataItem_ setDataItemEelements:responseObject];
        [self updateEditType];
        [self updateUIFromData];
    }
    else if ([[URLManager urlString:URLTypeOneClassUserCount variableKey:key] isEqualToString:candidate])
    {
        NSString* registeredCount = [responseObject objectForKey:[NotificationDataItem getKey:ElementType_ClassRegisteredCount ]];
        NSString* maxCount = [responseObject objectForKey:[NotificationDataItem getKey:ElementType_ClassStudent ]];
        int leftCount = (int)(maxCount.integerValue - registeredCount.integerValue);
        NSString* registeredInformation = [NSString stringWithFormat:@"Registered : %@, left : %i",registeredCount, leftCount];
        registerCountLabel_.text = registeredInformation;
        if  (editType_ == EditType_View)
        {
            registerButton_.hidden = leftCount <= 0;
        }
        else
        {
            registerButton_.hidden = YES;
        }
    }
    else if ([classPostURLString isEqualToString:candidate])
    {
        [self.navigationController popViewControllerAnimated:YES];
    }
    else if ([[URLManager urlString:URLTypeUserRegister] isEqualToString:candidate])
    {
        
    }
    else if ([URLManager isURLString:URLTypeQueryRegisterStatus candidate:candidate])
    {
        NSString* uniqueNameKey = [MeViewData singleton].userUniqueName;
        NSString* registerResult = [responseObject objectForKey:uniqueNameKey];
        if (registerResult != nil && [registerResult compare:@"YES"] == NSOrderedSame)
        {
            bCurrentUserRegistered_ = YES;
            [registerButton_ setTitle:@"UnRegister" forState:UIControlStateNormal];
        }
        else
        {
            bCurrentUserRegistered_ = NO;
            [registerButton_ setTitle:@"Register" forState:UIControlStateNormal];
        }
    }
    else if ([[URLManager urlString:URLTypeUserUnregister] isEqualToString:candidate])
    {
        // do nothing
    }
    
}



- (void)onFail:(AFHTTPRequestOperation *)operation error:(NSError*)error
{
    [super onFail:operation error:error];
    // Do nothing
}

- (void)viewWillAppear:(BOOL)animated
{
    bCurrentUserRegistered_ = YES;
    if (EditType_View == editType_
        || EditType_Editable == editType_)
    {
        [self getOneClass];
    }

    if (EditType_View == editType_)
    {
        [self queryRegisterStatus];
    }
    
    if (editType_ == EditType_View)
    {
        self.navigationItem.rightBarButtonItem.title = @"";
    }
    
    [self getOneClassUserCount];
    
    [self updateEditType];
    [self updateUIFromData];
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

- (void)setEditType:(EditType)type
{
    editType_ = type;
}

- (void)updateEditType
{
    BOOL editable = YES;
    BOOL deletable = YES;
    BOOL registable = NO;
    switch (editType_)
    {
        case EditType_View:
        {
            editable = NO;
            deletable = NO;
            registable = YES;
        }
        break;
        case EditType_Editable:
        {
            editable = YES;
        }
        break;
        case EditType_Add:
        {
            editable = YES;
            deletable = NO;
        }
        break;
            
        default:
            break;
    }
    [self setEditable:editable];
    
    deleteButton_.hidden = deletable == NO;
    registerButton_.hidden = registable == NO;
    self.navigationItem.rightBarButtonItem.enabled = editable;
}

- (void)setNotificationDataItem:(NotificationDataItem*)item
{
    notificationDataItem_ = item;
}

- (IBAction)onDelete:(id)sender
{
    [self deleteOneClass];
}


- (IBAction)onRegister:(id)sender
{
    if (bCurrentUserRegistered_)
    {
        [self unregisterOneUser];
    }
    else
    {
        [self registerOneUser];
    }
    
    [self.navigationController popViewControllerAnimated:YES];
}

- (IBAction)onStudentDetail:(id)sender
{
    [OneClassReigsterStudentData singleton].classUniqueID = [notificationDataItem_ getUniqueKey];
    NSString* title = [notificationDataItem_ getClassName];
    title = [title stringByAppendingString:@" Registered users"];
    oneClassRegisterStudentViewController_.title = title;
    [self.navigationController pushViewController:oneClassRegisterStudentViewController_ animated:YES];
}


- (void)setEditable:(BOOL)bEditable
{
    classNameTextView_.editable = bEditable;
    classTimeTextField_.enabled = bEditable;
    classTeachTextField_.enabled = bEditable;
    classDescriptionTextView_.editable = bEditable;
    classMaxStudentNumberTextField_.enabled = bEditable;
}

- (void)updateUIFromData
{
    classNameTextView_.text = [notificationDataItem_ getClassName];
    classTimeTextField_.text = [notificationDataItem_ getClassTime];
    classTeachTextField_.text = [notificationDataItem_ getClassTeacher];
    classDescriptionTextView_.text = [notificationDataItem_ getClassDescription];
    classMaxStudentNumberTextField_.text = [notificationDataItem_ getClassMaxStudent];
}


@end
